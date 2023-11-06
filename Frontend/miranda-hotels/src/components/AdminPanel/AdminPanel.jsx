import React, { useContext } from "react";
import { useNavigate } from 'react-router-dom';
import { UserProfileContext } from '../../App';
import "./AdminPanel.css"

function AdminPanel() {
    const { userProfile } = useContext(UserProfileContext);
    const navigate = useNavigate();

    if (!userProfile || userProfile.role !== "Admin") {
        return null;
      }

      return (
        <div className="contenedorAdmin">
            <div>
                <h4>Panel de Administración</h4>
            </div>
            <div>
                <button className="button" onClick={() => navigate('/loadhotel')}>Nuevo Hotel</button>
                <button className="button" onClick={() => navigate('/loadamenity')}>Nuevo Amenity</button>
                <button className="button" onClick={() => navigate('/admin/reservations/hotel')}>Ver reservas por Hotel</button>
                <button className="button" onClick={() => navigate('/admin/reservations/user')}>Ver reservas por usuario</button>
            </div>
        </div>
      )

}

export default AdminPanel