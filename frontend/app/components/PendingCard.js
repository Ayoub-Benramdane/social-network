export default function PendingGroupRequestCard({ group, onCancelRequest }) {
  return (
    <div className="pending-group-card">
      <div className="pending-group-content">
        <div className="pending-group-header">
          <div className="pending-group-info">
            <img
              src={group.avatar}
              alt={group.name}
              className="pending-group-avatar"
            />

            <div className="pending-group-name">
              <h3>{group.name}</h3>
              <div className="pending-group-status">
                <span>Request Pending</span>
              </div>
            </div>
          </div>

          <button
            onClick={() => onCancelRequest(group.id)}
            className="pending-group-cancel-btn"
          >
            Cancel
          </button>
        </div>

        {group.description && (
          <p className="pending-group-description">{group.description}</p>
        )}

        <div className="pending-group-request-date">
          Requested on {group.created_at || "2 hours ago"}
        </div>
      </div>
    </div>
  );
}
