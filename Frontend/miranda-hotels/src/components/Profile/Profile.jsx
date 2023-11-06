import React, {useContext, useEffect, useState} from "react";
import { useNavigate } from 'react-router-dom';
import Navbar from "../NavBar/NavBar";
import { UserProfileContext } from '../../App';
import AdminPanel from "../AdminPanel/AdminPanel";
import "./Profile.css"

function Profile() {
    const { userProfile, setUserProfile } = useContext(UserProfileContext);
    const [userData, setUserData] = useState({});
    const { error, setError } = useState(null)
    const navigate = useNavigate();

    useEffect(() => {
        const fetchUserData = async () => {
            if (userProfile) {
                try {
                    const response = await fetch(`http://localhost:8090/user/${userProfile.id}`);
                    if (response.ok) {
                        const data = await response.json();
                        setUserData(data);
                    } else {
                        const errorData = await response.json();
                        throw new Error(errorData.error);
                    }
                } catch (error) {
                    setError(error.message);
                }
            }
        };
        fetchUserData();
    }, [userProfile]);


    const handleLogout = () => {
        localStorage.removeItem('token');
        setUserProfile(null);
        navigate('/');
    };

    if (!userProfile) {
        return (
            <>
                <Navbar />
                <div className="descripcion">
                    <p>No puedes acceder a este sitio.</p>
                </div>
            </>
        )
      }

    if (error) {
        return (
            <>
                <Navbar />
                <div className="fullscreen">Error: {error}</div>
            </>
        );
    }

    if (userData) {
        return (
            <>
                <Navbar />
                <div className="descripcion">
                    <h3>Perfil de Usuario</h3>
                    <p>Nombre: {userData.name}</p>
                    <p>Apellido: {userData.last_name}</p>
                    <p>DNI: {userData.dni}</p>
                    <p>Email: {userData.email}</p>
                    <p>Nº de Usuario: {userData.id}</p>
                    <div>
                        {userProfile.role === "Customer" && <button className="button" onClick={()=>navigate(`/user/reservations/${userProfile.id}`)}> Mis Reservas </button>}
                        <button className="button" onClick={handleLogout}>Cerrar Sesión</button>
                    </div>
                    <AdminPanel />
                </div>
            </>
        )
    }
}

export default Profile;
