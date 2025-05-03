"use client";
import { useState, useEffect } from "react";
import "../../styles/GroupsPage.css";
import "../styles/PostComponent.css";
import Navbar from "../components/NavBar";
import EventCard from "../components/EventCard";

export default function GroupEvents({ groupId }) {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const currentUser = {
    first_name: "Mohammed Amine",
    last_name: "Dinani",
    avatar: "./avatars/thorfinn-vinland-saga-episode-23-1.png",
    username: "mdinani",
  };

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const response = await fetch("http://localhost:8404/events", {
          method: "GET",
          credentials: "include",
        });
  
        const data = await response.json();
        console.log("Fetched events aaaaa:", data);
  
        if (!response.ok) {
          throw new Error(data.error || "Failed to fetch events.");
        }
  
        const filteredEvents = data.filter((e) => e.group_id === groupId);
        setEvents(filteredEvents);
      } catch (err) {
        console.error("Error fetching events:", err);
        setError("An error occurred while fetching events.");
      } finally {
        setLoading(false);
      }
    };
  
    if (groupId) {
      fetchEvents();
    }
  }, [groupId]);

  return (
    <div className="group-events">
      <Navbar user={currentUser} />
      <h2>Group Events</h2>

      {loading && <p>Loading events...</p>}
      {error && <p className="error">{error}</p>}

      <div className="events-list">
        {events.length > 0 ? (
          events.map((event) => (
            <EventCard key={event.id} event={event} />
          ))
        ) : (
          !loading && <p>No events found for this group.</p>
        )}
      </div>
    </div>
  );
}
