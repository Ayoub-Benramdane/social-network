"use client";
import "../styles/SuggestedCommunities.css";

export default function SuggestedCommunities() {
  const communities = [
    { id: 1, name: "Volkswagen", avatar: "https://via.placeholder.com/40" },
    { id: 2, name: "Coffee", avatar: "https://via.placeholder.com/40" },
    { id: 3, name: "Zone01Oujda", avatar: "https://via.placeholder.com/40" },
    { id: 4, name: "Zone01Oujda", avatar: "https://via.placeholder.com/40" },
  ];

  return (
    <div className="suggested-communities">
      <h3>Community</h3>
      <ul className="community-list">
        {communities.map((community) => (
          <li key={community.id} className="community-item">
            <img
              src={community.avatar}
              alt={community.name}
              className="avatar"
            />
            <span className="community-name">{community.name}</span>
          </li>
        ))}
      </ul>
      <button className="view-more-btn">View more</button>
    </div>
  );
}
