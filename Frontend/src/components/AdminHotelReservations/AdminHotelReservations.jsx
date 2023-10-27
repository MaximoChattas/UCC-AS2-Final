import React, { useEffect, useState, useContext } from "react";
import { LoginContext, UserProfileContext } from '../../App';
import { Link } from "react-router-dom";
import Navbar from "../NavBar/NavBar";
import "./AdminHotelReservations.css"

const AdminHotelReservations = () => {
  const [hotelReservations, setHotelReservations] = useState({ reservations: [] });
  const [hotels, setHotels] = useState([]);
  const [error, setError] = useState(null);
  const { userProfile } = useContext(UserProfileContext);
  const { loggedIn } = useContext(LoginContext);

  useEffect(() => {
    const fetchHotelReservations = async () => {
      try {
        const response = await fetch(`http://localhost:8085/reservation`);
        if (response.ok) {
          const data = await response.json();
          setHotelReservations({ reservations: data });

          const hotelResponse = await fetch(`http://localhost:8085/hotel`);
          if (hotelResponse.ok) {
            const hotelData = await hotelResponse.json();
            setHotels(hotelData);
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

    fetchHotelReservations();
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!hotelReservations) {
    return <div>Loading...</div>;
  }

  if (!loggedIn || userProfile.role !== "Admin") {
    return (
      <>
        <Navbar />
        <p className="fullscreen">No puedes acceder a este sitio.</p>
      </>
    );
  }

  return (
    <>
      <Navbar />
      <h2>Reservas</h2>

      <div className="containerReservations">
      <ul className="list-group">
        {hotels.map(hotel => {
          const filteredReservations = hotelReservations.reservations || [];
          const hotelReservationsFiltered = filteredReservations.filter(
            reservation => reservation.hotel_id === hotel.id
          );
          return (
            <li key={hotel.id} className="list-group-item list-group-item-dark">
              <Link to={`/hotel/${hotel.id}`}>
                <h3>{hotel.name}</h3>
              </Link>
              {hotelReservationsFiltered.length > 0 ? (
                <ul className="list-group">
                  {hotelReservationsFiltered.map(reservation => (
                    <li key={reservation.id} className="list-group-item">
                      <Link to={`/reservation/${reservation.id}`} >
                        <p>Nº Reserva: {reservation.id}</p>
                      </Link>
                      <p>Inicio: {reservation.start_date}</p>
                      <p>Fin: {reservation.end_date}</p>
                      <p>Costo: {reservation.amount}</p>
                      <Link to={`/user/${reservation.user_id}`}>
                        <p>Nº Usuario: {reservation.user_id}</p>
                      </Link>
                    </li>
                  ))}
                </ul>
              ) : (
                <p>Sin reservas.</p>
              )}
            </li>
          );
        })}
      </ul>
      </div>
    </>
  );
};

export default AdminHotelReservations;
