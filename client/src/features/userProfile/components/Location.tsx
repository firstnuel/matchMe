import { useState, useEffect } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";
import { type Point } from "../../../shared/types/user";
import { useMutation } from "@tanstack/react-query";
import { getLocationCity } from "../api/userProfile";

interface LocationProps {
  location: Point;
  setLocation: React.Dispatch<React.SetStateAction<Point>>;
}

const Location = ({ location, setLocation }: LocationProps) => {
  const [error, setError] = useState("");
  const mutation = useMutation({
    mutationFn: getLocationCity,
    onSuccess: (data) => {
      setLocation({
        latitude: data.lat,
        longitude: data.lon,
        city: data.city,
      });
    },
    onError: () => {
      setError("Could not fetch city name.");
    },
  });

  useEffect(() => {
    if (
      location.latitude !== 0 &&
      location.longitude !== 0 &&
      location.city === ""
    ) {
      mutation.mutate({ latitude: location.latitude, longitude: location.longitude });
    }
  }, [location.latitude, location.longitude, location.city, mutation]);

  const getLocation = () => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        async (position) => {
          try {
            const { latitude, longitude } = position.coords;
            mutation.mutate({ latitude, longitude });
            setError("");
          } catch (err) {
            console.error(err);
            setError("Could not fetch city name.");
          }
        },
        (err) => {
          setError("Unable to retrieve location. Please allow location access.");
          console.error(err);
        }
      );
    } else {
      setError("Geolocation is not supported by this browser.");
    }
  };

  return (
    <div className="location-section">
      <div>
        <Icon icon="mdi:location" className="location-icon" />
        Your location helps us show you matches nearby
      </div>
      <button className="location-btn" onClick={getLocation}>
        Update Location
      </button>
      <div className="current-location">
        Current: {location.city || "unknown"}{" "}
        {location.latitude && location.longitude
          ? `(${location.latitude.toFixed(4)}, ${location.longitude.toFixed(4)})`
          : ""}
      </div>
      {error && <div className="info error">{error}</div>}
    </div>
  );
};

export default Location;