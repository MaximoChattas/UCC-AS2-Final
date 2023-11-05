import React, { useContext } from "react";
import { useNavigate } from 'react-router-dom';
import Navbar from "../NavBar/NavBar";
import { LoginContext, UserProfileContext } from '../../App';
import AdminPanel from "../AdminPanel/AdminPanel";
import "./Profile.css"

function Profile() {
    const { loggedIn, setLoggedIn } = useContext(LoginContext);
    const { userProfile, setUserProfile } = useContext(UserProfileContext);
    const navigate = useNavigate();

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('userProfile');
        setLoggedIn(false);
        setUserProfile(null);
        navigate('/');
    };

    if (!loggedIn) {
        return (
            <>
                <Navbar />
                <div className="descripcion">
                    <p>No puedes acceder a este sitio.</p>
                </div>
            </>
        )
      }

    const reservationURL = "/user/reservations/"+userProfile.id;

    return (
        <>
            <Navbar />
            <div className="descripcion">
                <h3>Perfil de Usuario</h3>
                <p>Nombre: {userProfile.name}</p>
                <p>Apellido: {userProfile.last_name}</p>
                <p>DNI: {userProfile.dni}</p>
                <p>Email: {userProfile.email}</p>
                <p>Nº de Usuario: {userProfile.id}</p>
                <div>
                    {userProfile.role === "Customer" && <button className="button" onClick={()=>navigate(reservationURL)}> Mis Reservas </button>}
                    {userProfile.role === "Customer" && <button className="button" onClick={()=>navigate("/user/reservations/range")}> Reservas por Rango </button>}
                    <button className="button" onClick={handleLogout}>Cerrar Sesión</button>
                </div>
                <AdminPanel />
            </div>
        </>
    )
}

export default Profile;
