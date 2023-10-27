import React, { useEffect, useState, useContext } from "react";
import { LoginContext, UserProfileContext } from '../../App';
import { useParams } from "react-router-dom";
import { Link } from "react-router-dom";
import Navbar from "../NavBar/NavBar";
import "./UserReservations.css"

const UserReservations = () => {
  const { id } = useParams();
  const [userReservations, setUserReservations] = useState([]);
  const [hotels, setHotels] = useState([]);
  const [error, setError] = useState(null);
  const { userProfile } = useContext(UserProfileContext);
  const { loggedIn } = useContext(LoginContext);

  useEffect(() => {
    const fetchUserReservations = async () => {
      try {
        const response = await fetch(`http://localhost:8085/user/reservations/${id}`);
        if (response.ok) {
          const data = await response.json();
          setUserReservations(data);

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

    fetchUserReservations();
  }, [id]);

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!userReservations) {
    return <div>Loading...</div>;
  }

  if (!loggedIn || ((userReservations.user_id !== userProfile.id) && (userProfile.role !== "Admin"))) {
    return (
        <>
            <Navbar />
            <p className="fullscreen">No puedes acceder a este sitio.</p>
        </>
    )
   }

   return (
    <>
      <Navbar />

      <div className="fullscreen">
        <h2>Reservas de {userReservations.user_name} {userReservations.user_last_name}</h2>

        {userReservations.reservations ? (
          <ul className="list-group">
            {hotels.map(hotel => {
              const hotelReservations = userReservations.reservations.filter(
                reservation => reservation.hotel_id === hotel.id
              );
              if (hotelReservations.length > 0) {
                return (
                  <li key={hotel.id} className="list-group-item list-group-item-dark">
                    <Link to={`/hotel/${hotel.id}`}>
                      <h3>{hotel.name}</h3>
                    </Link>
                    <ul className="list-group">
                      {hotelReservations.map(reservation => (
                        <li key={reservation.id} className="list-group-item">
                          <Link to={`/reservation/${reservation.id}`}>
                            Nº Reserva: {reservation.id}
                          </Link>
                          <p>Inicio: {reservation.start_date}</p>
                          <p>Fin: {reservation.end_date}</p>
                          <p>Costo: {reservation.amount}</p>
                        </li>
                      ))}
                    </ul>
                  </li>

                );
              } else {
                return null;
              }
            })}
          </ul>
        ) : (
          <p>No tiene ninguna reserva.</p>
        )}
      </div>
    </>
  );  
  
};

export default UserReservations;
