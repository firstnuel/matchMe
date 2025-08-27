import { useGetLocationCity } from './useCurrentUser';
import { useEffect, useState } from "react";

type Coordinates = {
  latitude?: number;
  longitude?: number;
  city?: string;
};

export function useCityFromCoordinates(coordinates?: Coordinates | null) {
  const cityMutation = useGetLocationCity();
  const [locationDisplay, setLocationDisplay] = useState<string | null>(null);
  const [fetchingCity, setFetchingCity] = useState(false);
  const [retryCount, setRetryCount] = useState(0);

  useEffect(() => {
    if (!coordinates) {
      setLocationDisplay(null);
      setRetryCount(0);
      return;
    }

    if (coordinates.city && coordinates.city.trim()) {
      setLocationDisplay(coordinates.city);
      setRetryCount(0);
      return;
    }

    if (
      coordinates.latitude &&
      coordinates.longitude &&
      coordinates.latitude !== 0 &&
      coordinates.longitude !== 0
    ) {
      setFetchingCity(true);

      cityMutation.mutate(
        {
          latitude: coordinates.latitude,
          longitude: coordinates.longitude,
        },
        {
          onSuccess: (data) => {
            setLocationDisplay(data.city);
            setFetchingCity(false);
            setRetryCount(0);
          },
          onError: () => {
            if (retryCount < 1) {
              setRetryCount((prev) => prev + 1); // retry once
            } else {
              if (coordinates.latitude && coordinates.longitude) {
                setLocationDisplay(
                  `(${coordinates.latitude.toFixed(4)}, ${coordinates.longitude.toFixed(4)})`
                );
              }
              setFetchingCity(false);
            }
          },
        }
      );
    } else {
      setLocationDisplay(null);
      setRetryCount(0);
    }
    // include retryCount so we can trigger retry
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [coordinates, retryCount]);

  return { locationDisplay, fetchingCity, retry: () => setRetryCount((c) => c + 1) };
}
