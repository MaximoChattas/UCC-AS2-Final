import React, { useContext, useState } from "react";
import { Link } from "react-router-dom";
import { LoginContext, UserProfileContext } from '../../App';
import Navbar from "../NavBar/NavBar";
import Calendar from "../Calendar/Calendar";
import { format } from "date-fns";
import "./UserReservationsRange.css"

const ReservationsInRange = () => {
  const [reservations, setReservations] = useState([]);
  const [error, setError] = useState(null);
  const [selectedDates, setSelectedDates] = useState({
    startDate: new Date(),
    endDate: new Date(),
  });

  const { loggedIn } = useContext(LoginContext);
  const { userProfile } = useContext(UserProfileContext);

    if (!loggedIn) {
        return (
            <>
                <Navbar />
                <p className="fullscreen">No puedes acceder a este sitio.</p>
            </>
        )
    }

  const id = userProfile.id;

  const fetchReservations = async () => {
    try {
      const startDate = format(selectedDates.startDate, "dd-MM-yyyy");
      const endDate = format(selectedDates.endDate, "dd-MM-yyyy");
      const startTime = "00:00";
      const endTime = "23:59";
      const startDateTime = `${startDate}+${startTime}`;
      const endDateTime = `${endDate}+${endTime}`;
      const url = `http://localhost:8090/user/reservations/${id}/range?start_date=${startDateTime}&end_date=${endDateTime}`;
      const response = await fetch(url);
      if (response.ok) {
        const data = await response.json();
        setReservations(data);
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

  const handleSelectDates = (selectedRange) => {
    setSelectedDates(selectedRange);
  };

  const formattedStartDate = format(selectedDates.startDate, "dd/MM/yyyy");
  const formattedEndDate = format(selectedDates.endDate, "dd/MM/yyyy");

   if (!reservations) {
    return (
        <>
            <Navbar />
            <p className="fullscreen">No tienes reservas en el rango de fechas seleccionado.</p>
        </>
    )
   }

  if (reservations.length === 0) {
    return (
      <>
        <Navbar />
        <div className="fullscreen">
            <h2>Ver Reservas</h2>
            <p>
              Seleccione un rango de fechas en el calendario para ver sus reservas.
            </p>
            <Calendar onSelectDates={handleSelectDates} />
            <button onClick={fetchReservations} className="button">Ver</button>
        </div>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <div className="fullscreen">
      <h2>Reservas de {userProfile.name} {userProfile.last_name}</h2>
      <h5>{formattedStartDate} - {formattedEndDate}</h5>
          <ul>
            {reservations.map(reservation => (
                <li key={reservation.id} className="list-group-item list-group-item-dark">
                <Link to={`/reservation/${reservation.id}`}>
                    Nº Reserva: {reservation.id}
                </Link>
                <p>Inicio: {reservation.start_date}</p>
                <p>Fin: {reservation.end_date}</p>
                <p>Costo: {reservation.amount}</p>
                </li>
        ))}
            </ul>
      </div>

    </>
  );
};

export default ReservationsInRange;