import React, { useState } from "react";
import { Routes, Route } from 'react-router-dom'
import HotelList from "./components/HotelList/HotelList";
import HotelDetails from "./components/HotelDetails/HotelDetails";
import Login from "./components/Login/Login";
import Profile from "./components/Profile/Profile";
import Signup from "./components/SignUp/SignUp";
import ReservationDetails from "./components/ReservationDetails/ReservationDetails";
import UserReservations from "./components/UserReservations/UserReservations";
import LoadHotel from "./components/LoadHotel/LoadHotel";
import AdminHotelReservations from "./components/AdminHotelReservations/AdminHotelReservations";
import AdminUserReservations from "./components/AdminUserReservations/AdminUserReservations";
import UserDetails from "./components/UserDetails/UserDetails";
import "./App.css"
import HotelAvailable from "./components/HotelAvailable/HotelAvailable";
import LoadAmenity from "./components/LoadAmenity/LoadAmenity.jsx";
import UpdateHotel from "./components/UpdateHotel/UpdateHotel.jsx";
import ContainerStats from "./components/ContainersStats/ContainerStats.jsx"
import ScaleServices from "./components/ScaleServices/ScaleServices.jsx";


// Create the UserProfileContext
export const UserProfileContext = React.createContext();

const App = () => {
  const [userProfile, setUserProfile] = useState(null);

  return (
    <div className="menu">
        <UserProfileContext.Provider value={{ userProfile, setUserProfile }}>
          <Routes>
            <Route path="/" exact element={<HotelList />} />
            <Route path="/hotel/:id" element={<HotelDetails />} />
            <Route path="/user/:id" element={<UserDetails />} />
            <Route path="/login" element={<Login />} />
            <Route path="/profile" element={<Profile />} />
            <Route path="/signup" element={<Signup />} />
            <Route path="/loadhotel" element={<LoadHotel />} />
            <Route path="/reservation/:id" element={<ReservationDetails />} />
            <Route path="/user/reservations/:id" element={<UserReservations />} />
            <Route path="/admin/reservations/hotel" element={<AdminHotelReservations />} />
            <Route path="/admin/reservations/user" element={<AdminUserReservations />} />
            <Route path="/hotel/availability" element={<HotelAvailable />} />
            <Route path="/loadamenity" element={<LoadAmenity />} />
            <Route path="/updatehotel/:id" element={<UpdateHotel />} />
            <Route path="/stats" element={<ContainerStats />} />
            <Route path="/scale" element={<ScaleServices />} />
          </Routes>
        </UserProfileContext.Provider>
    </div>
  );
};

export default App;
