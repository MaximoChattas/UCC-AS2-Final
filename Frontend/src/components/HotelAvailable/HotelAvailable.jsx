import React, { useContext, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { LoginContext } from '../../App';
import Navbar from "../NavBar/NavBar";
import "../HotelList/HotelList.css";
import Calendar from "../Calendar/Calendar";
import { format } from "date-fns";

const HotelAvailable = () => {
  const [hotels, setHotels] = useState([]);
  const [error, setError] = useState(null);
  const [selectedDates, setSelectedDates] = useState({
    startDate: new Date(),
    endDate: new Date(),
  });

  const { loggedIn } = useContext(LoginContext);

  const fetchHotels = async () => {
    try {
      const startDate = format(selectedDates.startDate, "dd-MM-yyyy");
      const endDate = format(selectedDates.endDate, "dd-MM-yyyy");
      const startTime = "15:00";
      const endTime = "11:00";
      const startDateTime = `${startDate}+${startTime}`;
      const endDateTime = `${endDate}+${endTime}`;
      const url = `http://localhost:8090/availability?start_date=${startDateTime}&end_date=${endDateTime}`;
      const response = await fetch(url);
      if (response.ok) {
        const data = await response.json();
        setHotels(data);
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

  if (!hotels) {
    return (
        <>
          <Navbar />
          <h2>Hoteles Disponibles</h2>
          <p className="fullscreen">No hay hoteles disponibles</p>
        </>
    );
  }

  if (hotels.length === 0) {
    return (
      <>
        <Navbar />
        <div className="fullscreen">
          <h2>Verificar Disponibilidad</h2>
          <p>
            Seleccione un rango de fechas en el calendario para verificar los hoteles
            disponibles
          </p>
          <Calendar onSelectDates={handleSelectDates} />
          <button onClick={fetchHotels} style={{ marginTop: '20px' }}>Verificar</button>
        </div>
      </>
    );
  }

  const formattedStartDate = format(selectedDates.startDate, "dd/MM/yyyy");
  const formattedEndDate = format(selectedDates.endDate, "dd/MM/yyyy");

  return (
    <>
      <Navbar />
      <h2>Hoteles Disponibles</h2>
      <h5>{formattedStartDate} - {formattedEndDate}</h5>
      <div className="row">
        {hotels.map((hotel) => (
          <div key={hotel.id} className="col-md-4 mb-4">
            <div className="card">
              {hotel.images &&
                  <img className="card-img-top"
                       alt={`Image for ${hotel.name}`}
                       src={`http://localhost:8090/image/${hotel.images[0].id}`}
                  />}
              <div className="card-body">
                <h5 className="card-title">
                  <Link to={`/hotel/${hotel.id}`}>
                    {hotel.name}
                  </Link>
                </h5>
                <p className="card-text">
                  Dirección: {hotel.street_name} {hotel.street_number}
                </p>
                <p className="card-text">${hotel.rate}</p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </>
  );
};

export default HotelAvailable;