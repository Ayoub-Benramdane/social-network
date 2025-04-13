export default function UserCard({ user, action, onClick }) {
  const handleFollow = (id) => {
    console.log("Followed user with ID:", id);
  };

  return (
    <li className="user-item" onClick={onClick}>
      <img
        src={user.avatar || user.image}
        className="user-avatar"
        alt={user.username || user.name}
      />
      <div className="user-details">
        <div className="user-info">
          <h4 className="user-name">
            {user.first_name
              ? `${user.first_name} ${user.last_name}`
              : user.name}
          </h4>
          <p className="user-username">
            {user.username
              ? `@${user.username}`
              : `(${user.total_members}) Members`}
          </p>
        </div>

        {action === "follow" && (
          <button className="follow-btn" onClick={() => handleFollow(user.id)}>
            Follow
          </button>
        )}
      </div>
    </li>
  );
}
