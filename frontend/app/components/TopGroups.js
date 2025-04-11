"use client";
import "../styles/TopGroups.css";

export default function TopGroups({ groups }) {
  return (
    <div className="suggested-groups">
      <div className="header">
        <h3>Groups</h3>
        <p style={{ color: "#3555F9", fontSize: "14px", cursor: "pointer" }}>
          See all &rarr;
        </p>
      </div>
      <ul className="group-list">
        {groups.map((group) => (
          <ul key={group.id} className="group-item">
            <img
              src={group.image}
              // alt={group.name}
              className="avatar"
            />
            <span className="group-name">{group.name}</span>
          </ul>
        ))}
      </ul>
    </div>
  );
}
