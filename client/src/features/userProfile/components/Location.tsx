import { useState, useEffect } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";
import { type Point } from "../../../shared/types/user";
import { useGetLocationCity } from "../hooks/useCurrentUser";

interface LocationProps {
  location: Point;
  setLocation: React.Dispatch<React.SetStateAction<Point>>;
  maxRetries?: number; // optional, defaults to 1 retry
}

const Location = ({ location, setLocation, maxRetries = 1 }: LocationProps) => {
  const [error, setError] = useState("");
  const [retryCount, setRetryCount] = useState(0);
  const mutation = useGetLocationCity();

  useEffect(() => {
    if (
      location.latitude !== 0 &&
      location.longitude !== 0 &&
      (!location.city || location.city === "")
    ) {
      mutation.mutate(
        { latitude: location.latitude, longitude: location.longitude },
        {
          onSuccess: (data) => {
            setLocation({
              latitude: data.lat,
              longitude: data.lon,
              city: data.city,
            });
            setError("");
            setRetryCount(0); // reset retries on success
          },
          onError: () => {
            if (retryCount < maxRetries) {
              setRetryCount((prev) => prev + 1); // retry
            } else {
              setError("Could not fetch city name.");
            }
          },
        }
      );
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location, retryCount]);

  const getLocation = () => {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        async (position) => {
          try {
            const { latitude, longitude } = position.coords;
            mutation.mutate(
              { latitude, longitude },
              {
                onSuccess: (data) => {
                  setLocation({
                    latitude: data.lat,
                    longitude: data.lon,
                    city: data.city,
                  });
                  setError("");
                  setRetryCount(0); // reset on success
                },
                onError: () => {
                  if (retryCount < maxRetries) {
                    setRetryCount((prev) => prev + 1);
                  } else {
                    setError("Could not fetch city name.");
                  }
                },
              }
            );
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
