import { useState } from "react";
import "../styles/GroupFormStyle.css";

export default function GroupForm() {
  const [groupFormInputs, setGroupFormInputs] = useState({
    name: "",
    description: "",
    groupImage: null,
    privacy: "",
    members: [],
  });
  const [imageInputKey, setImageInputKey] = useState(Date.now());

  const handleImageChange = (event) => {
    const file = event.target.files[0];
    setGroupFormInputs({
      ...groupFormInputs,
      groupImage: file,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const formData = new FormData();

    // formData.append(
    //   "groupData",
    //   JSON.stringify({
    //     name: groupFormInputs.name,
    //     description: groupFormInputs.description,
    //     members: groupFormInputs.members,
    //   })
    // );
    const fieldsToInclude = ["name", "description", "privacy"];

    fieldsToInclude.forEach((field) => {
      formData.append(field, groupFormInputs[field]);
    });

    if (groupFormInputs.groupImage) {
      formData.append("groupImage", groupFormInputs.groupImage);
    }

    for (const pair of formData.entries()) {
      console.log(pair);
    }
    try {
      const response = await fetch("http://localhost:8404/new_group", {
        method: "POST",
        credentials: "include",
        body: formData,
      });

      if (!response.ok) {
        const data = await response.json();
        console.log(data);

        throw new Error(data.error || "Failed to create the group");
      }

      console.log("Group created successfully");
      setGroupFormInputs({
        name: "",
        description: "",
        members: "",
        groupImage: null,
      });
    } catch (error) {
      console.log(error);
    }
  };

  return (
    <div className="groupDiv">
      <h3 style={{ color: "#18151b" }}>Create a new Group</h3>
      <form className="create-group-form" onSubmit={handleSubmit}>
        <div className="form-div">
          <div className="name-input">
            <input
              placeholder="Group name..."
              required
              value={groupFormInputs.name}
              onChange={(e) => {
                setGroupFormInputs({
                  ...groupFormInputs,
                  name: e.target.value,
                });
              }}
            ></input>
          </div>

          <div className="description-input">
            <input
              placeholder="Group description..."
              required
              value={groupFormInputs.description}
              onChange={(e) => {
                setGroupFormInputs({
                  ...groupFormInputs,
                  description: e.target.value,
                });
              }}
            ></input>
          </div>
          <div className="image-input">
            <p>Upload Image</p>
            {groupFormInputs.groupImage && (
              <div>
                <img
                  style={{
                    width: "150px",
                    borderRadius: "8px",
                  }}
                  src={URL.createObjectURL(groupFormInputs.groupImage)}
                  alt="Selected"
                />
                <button
                  onClick={(e) => {
                    e.preventDefault();
                    setGroupFormInputs({
                      ...groupFormInputs,
                      groupImage: null,
                    });
                    setImageInputKey(Date.now());
                  }}
                >
                  Remove
                </button>
              </div>
            )}
            <input
              key={imageInputKey}
              type="file"
              name="groupImage"
              onChange={handleImageChange}
            ></input>
          </div>
        </div>
        <button
          style={{
            backgroundColor: "rgb(174, 0, 255)",
            width: "150px",
            height: "32px",
            border: "none",
            color: "white",
            borderRadius: "8px",
            fontSize: "14px",
            fontWeight: "500",
          }}
        >
          Create group
        </button>
      </form>
    </div>
  );
}
