import React, { useEffect, useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { useParams } from "react-router-dom";
import { Link } from "react-router-dom";
import { UserProfileContext } from '../../App';
import "./ReservationDetail.css"

import Navbar from "../NavBar/NavBar";

const ReservationDetails = () => {
  const { id } = useParams();
  const [reservation, setReservation] = useState(null);
  const [hotel, setHotel] = useState(null);
  const [error, setError] = useState(null);
  const [deleteError, setDeleteError] = useState(null)
  const [limitDate, setLimitDate] = useState(new Date())

  const { userProfile } = useContext(UserProfileContext);
  const navigate = useNavigate()

  useEffect(() => {
    const fetchReservationDetails = async () => {
      try {
        const response = await fetch(`http://localhost:8090/reservation/${id}`);
        if (response.ok) {
          const data = await response.json();
          setReservation(data);

          const hotelResponse = await fetch(`http://localhost:8085/hotel/${data.hotel_id}`);
          if (hotelResponse.ok) {
            const hotelData = await hotelResponse.json();
            setHotel(hotelData);
          } else {
            const errorData = await hotelResponse.json();
            throw new Error(errorData.error);
          }

        } else {
          const errorData = await response.json();
          throw new Error(errorData.error);
        }
      } catch (error) {
        setError(error.message);
      }
    };

    fetchReservationDetails();
  }, [id]);

  const isReservationDeletable = () => {

    const [date, time] = reservation.start_date.split(' ');
    const [day, month, year] = date.split('-');
    const [hours, minutes] = time.split(':')

    const startDate = new Date(year, month - 1, day, hours, minutes);
    const currentDate = new Date()
    const differenceInHours = (startDate - currentDate) / 36e5;

    return differenceInHours >= 48;
  };

  const handleDeleteReservation = async () => {
    try {
      const response = await fetch(`http://localhost:8090/reservation/${id}`, {
        method: 'DELETE',
      });
      if (response.ok) {
        navigate(`/user/reservations/${userProfile.id}`)
      } else {
        const errorData = await response.json();
        throw new Error(errorData.error);
      }
    } catch (error) {
      setDeleteError(error.message);
    }
  };

  if (error) {
    return (
        <>
          <Navbar />
          <p className="fullscreen">{error}</p>
        </>
    );
  }

  if (!reservation) {
    return (
        <>
          <Navbar />
          <p className="fullscreen">Cargando...</p>
        </>
    );
  }

  if (!userProfile || ((reservation.user_id !== userProfile.id) && (userProfile.role === "Customer"))) {
    return (
        <>
          <Navbar />
          <p>No puedes acceder a este sitio.</p>
        </>
    )
  }

  return (
      <>
        <Navbar />
        <div className="Detalles">
          <h2>Reserva {id}</h2>
          {hotel && (
              <>
                <Link to={`/hotel/${hotel.id}`}>
                  <h4>{hotel.name}</h4>
                </Link>
                <h6>{hotel.street_name} {hotel.street_number}</h6>
              </>
          )}
          <p>Inicio: {reservation.start_date}</p>
          <p>Fin: {reservation.end_date}</p>
          <p>Costo: ${reservation.amount}</p>
          {isReservationDeletable() && userProfile.role === "Customer" &&
              <div>
                <button onClick={handleDeleteReservation}>Cancelar reserva</button>
                {deleteError && <p className="error-message">{deleteError}</p> }
              </div>}

        </div>
      </>
  );
};

export default ReservationDetails;
