export async function logout(){
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        await fetch(`http://localhost:8080/api/logout`, {
          method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',  
        })        
}