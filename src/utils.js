export const verifyToken = async (setIsLoggedIn) => {
    const token = localStorage.getItem('token');
    const res = await fetch(`${import.meta.env.VITE_API_URL}/verify`, {
        headers: { 'Authorization': `Bearer ${token}` }
    });
    if (res.status === 401) {
        localStorage.removeItem('token');
        setIsLoggedIn(false);
        return false;
    }
    setIsLoggedIn(true);
    return true;
};
