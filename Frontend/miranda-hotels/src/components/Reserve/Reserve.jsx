import React, { useContext, useState } from "react";
import { useNavigate } from "react-router-dom";
import { UserProfileContext } from '../../App';
import { format, differenceInHours } from "date-fns";

const Reservation = ({ hotel_id, hotelRate, startDate, endDate }) => {
  const { userProfile } = useContext(UserProfileContext);
  const [error, setError] = useState("");
  const [reservationSaved, setReservationSaved] = useState(false);
  const navigate = useNavigate();

  const calculateAmount = () => {
    const checkInDate = new Date(startDate);
    checkInDate.setHours(15, 0, 0, 0);

    const checkoutDate = new Date(endDate);
    checkoutDate.setHours(11, 0, 0, 0);

    const hours = differenceInHours(checkoutDate, checkInDate);
    const nights = Math.ceil(hours / 24);
    let amount = nights * hotelRate;
    amount = amount.toFixed(2);
    amount = parseFloat(amount).toLocaleString();

    return amount;
  };

  const handleReservation = async () => {
    const checkInDate = new Date(startDate);
    checkInDate.setHours(15, 0, 0, 0);

    const checkoutDate = new Date(endDate);
    checkoutDate.setHours(11, 0, 0, 0);

    try {
      const reservationData = {
        user_id: parseInt(userProfile.id),
        hotel_id: parseInt(hotel_id),
        start_date: format(checkInDate, "dd-MM-yyyy HH:mm"),
        end_date: format(checkoutDate, "dd-MM-yyyy HH:mm"),
      };

      const response = await fetch("http://localhost:8090/reserve", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(reservationData),
      });

      if (response.ok) {
        setReservationSaved(true);
        const data = await response.json();
        const url = "/reservation/" + data.id;
        navigate(url);
      } else {
        const data = await response.json();
        const errorMessage = data.error || "Error";
        throw new Error(errorMessage);
      }
    } catch (error) {
      console.error(error);
      setError(error.message);
    }
  };

  const amount = calculateAmount();

  if (userProfile.role === "Admin") {
    return null;
  }

  return (
    <div>
      <p>Total: ${amount}</p>
      <button onClick={handleReservation} disabled={reservationSaved}>
        {reservationSaved ? "Reservado" : "Reservar"}
      </button>
      {error && <p className="error-message">{error}</p>}
    </div>
  );
};

export default Reservation;