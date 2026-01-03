package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	structs "social-network/data"
	"social-network/database"
)

func InvitationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckLastActionTime(w, r, "invitations") {
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	type InvitationRequest struct {
		TargetUserID  int64 `json:"user_id"`
		TargetGroupID int64 `json:"group_id"`
	}

	var invitationRequest InvitationRequest
	err = json.NewDecoder(r.Body).Decode(&invitationRequest)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	invitationExists, err := database.InvitationExists(currentUser.UserID, invitationRequest.TargetUserID, invitationRequest.TargetGroupID)
	if err != nil {
		fmt.Println("Failed to retrieve invitation", err)
		response := map[string]string{"error": "Failed to retrieve invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	var existingInvitationID int64
	if invitationExists {
		existingInvitationID, err = database.GetInvitationID(currentUser.UserID, invitationRequest.TargetUserID, invitationRequest.TargetGroupID)
		if err != nil {
			fmt.Println("Failed to retrieve invitation ID", err)
			response := map[string]string{"error": "Failed to retrieve invitation ID"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	if invitationRequest.TargetGroupID == 0 {
		targetUser, err := database.FindUserByID(invitationRequest.TargetUserID)
		if err != nil {
			fmt.Println("Failed to retrieve follower", err)
			response := map[string]string{"error": "Failed to retrieve follower"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		} else if targetUser.UserID == currentUser.UserID {
			fmt.Println("Cannot send invitation to you", err)
			response := map[string]string{"error": "Cannot send invitation to you"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		isAlreadyFollowing, err := database.IsUserFollowing(currentUser.UserID, invitationRequest.TargetUserID)
		if err != nil {
			fmt.Println("Failed to check if user is followed", err)
			response := map[string]string{"error": "Failed to check if user is followed"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		var actionResult string
		if !isAlreadyFollowing {
			if targetUser.PrivacyLevel == "public" {
				if err := database.FollowUser(currentUser.UserID, invitationRequest.TargetUserID); err != nil {
					fmt.Println("Failed to follow user", err)
					response := map[string]string{"error": "Failed to follow user"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				if err := database.CreateNotification(currentUser.UserID, invitationRequest.TargetUserID, 0, 0, 0, "follow"); err != nil {
					fmt.Println("Failed to create notification", err)
					response := map[string]string{"error": "Failed to create notification"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}
				actionResult = "Unfollow"
			} else if !invitationExists {
				if err := database.CreateInvitation(currentUser.UserID, invitationRequest.TargetUserID, invitationRequest.TargetGroupID); err != nil {
					fmt.Println("Failed to send invitation", err)
					response := map[string]string{"error": "Failed to send invitation"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				if err := database.CreateNotification(currentUser.UserID, invitationRequest.TargetUserID, 0, 0, 0, "follow_request"); err != nil {
					fmt.Println("Failed to create notification", err)
					response := map[string]string{"error": "Failed to create notification"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				actionResult = "Pending"
			} else {
				if err := database.DeleteInvitation(existingInvitationID); err != nil {
					fmt.Println("Failed to delete invitation", err)
					response := map[string]string{"error": "Failed to delete invitation"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				if err := database.DeleteNotification(currentUser.UserID, invitationRequest.TargetUserID, 0, 0, 0, "follow_request"); err != nil {
					fmt.Println("Failed to create notification", err)
					response := map[string]string{"error": "Failed to delete notification"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				actionResult = "Follow"
			}
		} else {
			if err := database.UnfollowUser(currentUser.UserID, invitationRequest.TargetUserID); err != nil {
				fmt.Println("Failed to unfollow user", err)
				response := map[string]string{"error": "Failed to unfollow user"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			if err := database.DeleteNotification(currentUser.UserID, invitationRequest.TargetUserID, 0, 0, 0, "follow"); err != nil {
				fmt.Println("Failed to create notification", err)
				response := map[string]string{"error": "Failed to delete notification"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			actionResult = "Follow"
		}
		
		totalFollowingCount, err := database.CountUserFollowing(currentUser.UserID)
		if err != nil {
			fmt.Println("Failed to retrieve total following", err)
			response := map[string]string{"error": "Failed to retrieve total follows"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		targetUserFollowsBack, err := database.IsUserFollowing(invitationRequest.TargetUserID, currentUser.UserID)
		if err != nil {
			fmt.Println("Failed to check if user is followed", err)
			response := map[string]string{"error": "Failed to check if user is followed"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if actionResult == "Follow" && targetUserFollowsBack {
			actionResult = "Follow back"
		}

		responseData := struct {
			Action         string `json:"action"`
			TotalFollowers int64  `json:"total_followers"`
		}{
			Action:         actionResult,
			TotalFollowers: totalFollowingCount,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	} else {
		targetGroup, err := database.GetGroupByID(invitationRequest.TargetGroupID)
		if err != nil {
			fmt.Println("Failed to retrieve group", err)
			response := map[string]string{"error": "Failed to retrieve group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		isAlreadyMember, err := database.IsUserGroupMember(currentUser.UserID, invitationRequest.TargetGroupID)
		if err != nil {
			fmt.Println("Failed to check if user is member", err)
			response := map[string]string{"error": "Failed to check if user is member"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		if !isAlreadyMember {
			if targetGroup.PrivacyLevel == "public" {
				if err := database.AddUserToGroup(currentUser.UserID, invitationRequest.TargetGroupID); err != nil {
					fmt.Println("Failed to join group", err)
					response := map[string]string{"error": "Failed to join group"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}

				if err := database.CreateNotification(currentUser.UserID, targetGroup.AdminID, 0, targetGroup.GroupID, 0, "join"); err != nil {
					fmt.Println("Failed to create notification", err)
					response := map[string]string{"error": "Failed to create notification"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("join")
			} else if !invitationExists {
				if err := database.CreateInvitation(currentUser.UserID, targetGroup.AdminID, targetGroup.GroupID); err != nil {
					fmt.Println("Failed to send invitation", err)
					response := map[string]string{"error": "Failed to send invitation"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}
				if err := database.CreateNotification(currentUser.UserID, targetGroup.AdminID, 0, targetGroup.GroupID, 0, "join_request"); err != nil {
					fmt.Println("Failed to create notification", err)
					response := map[string]string{"error": "Failed to create notification"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("cancel")
			} else {
				if err := database.DeleteInvitation(existingInvitationID); err != nil {
					fmt.Println("Failed to delete invitation", err)
					response := map[string]string{"error": "Failed to delete invitation"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}
				if err := database.DeleteNotification(currentUser.UserID, targetGroup.AdminID, targetGroup.GroupID, 0, 0, "join_request"); err != nil {
					fmt.Println("Failed to create notification", err)
					response := map[string]string{"error": "Failed to delete notification"}
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(response)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode("follow")
			}
		} else if targetGroup.AdminID != currentUser.UserID {
			if err := database.RemoveUserFromGroup(currentUser.UserID, invitationRequest.TargetGroupID); err != nil {
				fmt.Println("Failed to leave group", err)
				response := map[string]string{"error": "Failed to leave group"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			if err := database.DeleteNotification(currentUser.UserID, targetGroup.AdminID, targetGroup.GroupID, 0, 0, "join"); err != nil {
				fmt.Println("Failed to create notification", err)
				response := map[string]string{"error": "Failed to delete notification"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("leave")
		} else {
			if err := database.DeleteGroup(invitationRequest.TargetGroupID); err != nil {
				fmt.Println("Failed to delete group", err)
				response := map[string]string{"error": "Failed to delete group"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("delete")
		}
	}
}

func GetGroupInvitations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	groupInvitations, err := database.GetGroupInvitations(currentUser.UserID)
	if err != nil {
		fmt.Println("Failed to retrieve invitations", err)
		response := map[string]string{"error": "Failed to retrieve invitations"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groupInvitations)
}

func AcceptInvitationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckLastActionTime(w, r, "follows") {
		return
	} else if !CheckLastActionTime(w, r, "group_members") {
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	type InvitationRequest struct {
		SenderUserID  int64 `json:"user_id"`
		RelatedGroupID int64 `json:"group_id"`
	}

	var invitationRequest InvitationRequest
	err = json.NewDecoder(r.Body).Decode(&invitationRequest)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Println("hadi l3adia", invitationRequest)

	var senderUser structs.User
	if invitationRequest.SenderUserID != 0 {
		senderUser, err = database.FindUserByID(invitationRequest.SenderUserID)
		if err != nil {
			fmt.Println("Failed to retrieve follower", err)
			response := map[string]string{"error": "Failed to retrieve follower"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		} else if senderUser.UserID == currentUser.UserID {
			fmt.Println("Cannot accept your own invitation", err)
			response := map[string]string{"error": "Cannot accept your own invitation"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		senderAlreadyFollows, err := database.IsUserFollowing(invitationRequest.SenderUserID, currentUser.UserID)
		if err != nil {
			fmt.Println("Failed to check if user is followed", err)
			response := map[string]string{"error": "Failed to check if user is followed"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		} else if senderAlreadyFollows && invitationRequest.RelatedGroupID == 0 {
			fmt.Println("User is already followed")
			response := map[string]string{"error": "User is already followed"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	if invitationRequest.RelatedGroupID != 0 {
		relatedGroup, err := database.GetGroupByID(invitationRequest.RelatedGroupID)
		if err != nil {
			fmt.Println("Failed to retrieve group", err)
			response := map[string]string{"error": "Failed to retrieve group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if relatedGroup.AdminID != invitationRequest.SenderUserID {
			fmt.Println("User is not the creator of the groupjhg")
			response := map[string]string{"error": "User is not the creator of the group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		currentUserIsMember, err := database.IsUserGroupMember(currentUser.UserID, relatedGroup.GroupID)
		if err != nil {
			fmt.Println("Failed to check if user is member of the group", err)
			response := map[string]string{"error": "Failed to check if user is member of the group"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		} else if currentUserIsMember {
			fmt.Println("User is already a member of the group")
			response := map[string]string{"error": "User is already a member of the group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	invitationExists, err := database.InvitationExists(invitationRequest.SenderUserID, currentUser.UserID, invitationRequest.RelatedGroupID)
	if err != nil || !invitationExists {
		fmt.Println("Failed to retrieve invitation", err)
		response := map[string]string{"error": "Failed to retrieve invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	invitationID, err := database.GetInvitationID(invitationRequest.SenderUserID, currentUser.UserID, invitationRequest.RelatedGroupID)
	if err != nil {
		fmt.Println("Failed to retrieve invitation ID", err)
		response := map[string]string{"error": "Failed to retrieve invitation ID"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.AcceptInvitation(invitationID, invitationRequest.SenderUserID, currentUser.UserID, invitationRequest.RelatedGroupID); err != nil {
		fmt.Println("Failed to accept invitation", err)
		response := map[string]string{"error": "Failed to accept invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func DeclineInvitationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckLastActionTime(w, r, "follows") {
		return
	} else if !CheckLastActionTime(w, r, "group_members") {
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	type InvitationRequest struct {
		SenderUserID  int64 `json:"user_id"`
		RelatedGroupID int64 `json:"group_id"`
	}

	var invitationRequest InvitationRequest
	err = json.NewDecoder(r.Body).Decode(&invitationRequest)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	senderUser, err := database.FindUserByID(invitationRequest.SenderUserID)
	if err != nil {
		fmt.Println("Failed to retrieve follower", err)
		response := map[string]string{"error": "Failed to retrieve follower"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else if senderUser.UserID == currentUser.UserID {
		fmt.Println("Cannot accept your own invitation", err)
		response := map[string]string{"error": "Cannot accept your own invitation"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	currentUserFollowsSender, err := database.IsUserFollowing(currentUser.UserID, invitationRequest.SenderUserID)
	if err != nil {
		fmt.Println("Failed to check if user is followed", err)
		response := map[string]string{"error": "Failed to check if user is followed"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	} else if currentUserFollowsSender && invitationRequest.RelatedGroupID == 0 {
		fmt.Println("User is already followed")
		response := map[string]string{"error": "User is already followed"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if invitationRequest.RelatedGroupID != 0 {
		_, err := database.GetGroupByID(invitationRequest.RelatedGroupID)
		if err != nil {
			fmt.Println("Failed to retrieve group", err)
			response := map[string]string{"error": "Failed to retrieve group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	invitationExists, err := database.InvitationExists(invitationRequest.SenderUserID, currentUser.UserID, invitationRequest.RelatedGroupID)
	if err != nil || !invitationExists {
		fmt.Println("Failed to retrieve invitation", err)
		response := map[string]string{"error": "Failed to retrieve invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	invitationID, err := database.GetInvitationID(invitationRequest.SenderUserID, currentUser.UserID, invitationRequest.RelatedGroupID)
	if err != nil {
		fmt.Println("Failed to retrieve invitation ID", err)
		response := map[string]string{"error": "Failed to retrieve invitation ID"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := database.DeleteInvitation(invitationID); err != nil {
		fmt.Println("Failed to decline invitation", err)
		response := map[string]string{"error": "Failed to decline invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}

func AcceptOtherInvitationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Method not allowed", r.Method)
		response := map[string]string{"error": "Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	if !CheckLastActionTime(w, r, "group_members") {
		return
	}

	currentUser, err := GetUserFromSession(r)
	if err != nil || currentUser == nil {
		fmt.Println("Failed to retrieve user", err)
		response := map[string]string{"error": "Failed to retrieve user"}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	type InvitationRequest struct {
		InvitedUserID  int64 `json:"user_id"`
		RelatedGroupID int64 `json:"group_id"`
		IsOwnerAction  bool  `json:"owner"`
	}

	var invitationRequest InvitationRequest
	err = json.NewDecoder(r.Body).Decode(&invitationRequest)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		response := map[string]string{"error": "Invalid request body"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println("hadi l other", invitationRequest)

	_, err = database.FindUserByID(invitationRequest.InvitedUserID)
	if err != nil {
		fmt.Println("Failed to retrieve follower", err)
		response := map[string]string{"error": "Failed to retrieve follower"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	} else if invitationRequest.InvitedUserID == currentUser.UserID {
		fmt.Println("Cannot accept your own invitation", err)
		response := map[string]string{"error": "Cannot accept your own invitation"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	relatedGroup, err := database.GetGroupByID(invitationRequest.RelatedGroupID)
	if err != nil {
		fmt.Println("Failed to retrieve group", err)
		response := map[string]string{"error": "Failed to retrieve group"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	invitedUserIsMember, err := database.IsUserGroupMember(invitationRequest.InvitedUserID, relatedGroup.GroupID)
	if err != nil {
		fmt.Println("Failed to check if user is member of the group", err)
		response := map[string]string{"error": "Failed to check if user is member of the group"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	invitationExists, err := database.InvitationExists(invitationRequest.InvitedUserID, currentUser.UserID, invitationRequest.RelatedGroupID)
	if err != nil || !invitationExists {
		fmt.Println("Failed to retrieve invitation", err)
		response := map[string]string{"error": "Failed to retrieve invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	invitationID, err := database.GetInvitationID(invitationRequest.InvitedUserID, currentUser.UserID, invitationRequest.RelatedGroupID)
	if err != nil {
		fmt.Println("Failed to retrieve invitation ID", err)
		response := map[string]string{"error": "Failed to retrieve invitation ID"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if invitationRequest.IsOwnerAction {
		if invitedUserIsMember {
			fmt.Println("User is already a member of the group")
			response := map[string]string{"error": "User is already a member of the group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if relatedGroup.AdminID != currentUser.UserID {
			fmt.Println("You are not the admin of the group")
			response := map[string]string{"error": "You are not the admin of the group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := database.AddUserToGroup(invitationRequest.InvitedUserID, invitationRequest.RelatedGroupID); err != nil {
			fmt.Println("Failed to join group", err)
			response := map[string]string{"error": "Failed to join group"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		
		if err := database.CreateNotification(invitationRequest.InvitedUserID, relatedGroup.AdminID, 0, relatedGroup.GroupID, 0, "join"); err != nil {
			fmt.Println("Failed to create notification", err)
			response := map[string]string{"error": "Failed to create notification"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		if !invitedUserIsMember {
			fmt.Println("User is not a member of the group")
			response := map[string]string{"error": "User is not a member of the group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		currentUserIsMember, err := database.IsUserGroupMember(currentUser.UserID, invitationRequest.RelatedGroupID)
		if err != nil {
			fmt.Println("Failed to check if user is member of the group", err)
			response := map[string]string{"error": "Failed to check if user is member of the group"}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		} else if currentUserIsMember {
			fmt.Println("User is already a member of the group")
			response := map[string]string{"error": "User is already a member of the group"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		if relatedGroup.PrivacyLevel == "public" {
			if err := database.AddUserToGroup(currentUser.UserID, invitationRequest.RelatedGroupID); err != nil {
				fmt.Println("Failed to join group", err)
				response := map[string]string{"error": "Failed to join group"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			if err := database.CreateNotification(currentUser.UserID, relatedGroup.AdminID, 0, relatedGroup.GroupID, 0, "join"); err != nil {
				fmt.Println("Failed to create notification", err)
				response := map[string]string{"error": "Failed to create notification"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		} else {
			if err := database.CreateInvitation(currentUser.UserID, relatedGroup.AdminID, invitationRequest.RelatedGroupID); err != nil {
				fmt.Println("Failed to send invitation", err)
				response := map[string]string{"error": "Failed to send invitation"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			if err := database.CreateNotification(currentUser.UserID, relatedGroup.AdminID, 0, relatedGroup.GroupID, 0, "join_request"); err != nil {
				fmt.Println("Failed to create notification", err)
				response := map[string]string{"error": "Failed to create notification"}
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}
	
	if err := database.DeleteInvitation(invitationID); err != nil {
		fmt.Println("Failed to delete invitation", err)
		response := map[string]string{"error": "Failed to delete invitation"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}