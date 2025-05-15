async function joinGroup(group) {
  try {
    const response = await fetch(
      `http://localhost:8404/groups/${group.group_id}`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      }
    );
  } catch (err) {
    console.log(err);
  }
}
async function leaveGroup(group) {
  try {
    const response = await fetch(
      `http://localhost:8404/groups/${group.group_id}`,
      {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      }
    );
  } catch (err) {
    console.log(err);
  }
}
async function deleteGroup(group) {
  try {
    const response = await fetch(
      `http://localhost:8404/groups/${group.group_id}`,
      {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
      }
    );
  } catch (err) {
    console.log(err);
  }
}

export { joinGroup, leaveGroup, deleteGroup };
