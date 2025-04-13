"use client";
import "../styles/TopGroups.css";

export default function TopGroups({ groups }) {
  return (
    <div className="top-groups">
      <div className="section-header">
        <h3>Suggested Groups</h3>
        <button className="see-all-btn">
          See all <span className="arrow">→</span>
        </button>
      </div>

      <ul className="group-list">
        {groups.map((group) => (
          <li key={group.id} className="group-item">
            <div className="group-image">
              <img
                src={group.image || "/images/default-group.png"}
                alt={group.name}
              />
            </div>
            <div className="group-info">
              <span className="group-name">{group.name}</span>
              <span className="group-members">
                {group.total_members || 0} members
              </span>
            </div>
            <button className="join-button">Join</button>
          </li>
        ))}
      </ul>

      {groups.length === 0 && (
        <div className="empty-state">
          <p>No groups available</p>
        </div>
      )}
    </div>
  );
}
