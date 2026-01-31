import React, { useState, useEffect } from 'react';

function App() {
  const [employees, setEmployees] = useState([]);
  const [customers, setCustomers] = useState([]);

  useEffect(() => {
    // Mengambil data dari Service Employee (Docker Port 8080)
    fetch('http://localhost:8080/employees')
      .then(res => res.json())
      .then(data => setEmployees(data))
      .catch(err => console.error("Employee API Error:", err));

    // Mengambil data dari Service Customer (Docker Port 8081)
    fetch('http://localhost:8081/customers')
      .then(res => res.json())
      .then(data => setCustomers(data))
      .catch(err => console.error("Customer API Error:", err));
  }, []);

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial' }}>
      <h1>IDStar Dashboard - Edwin Iqbal Santoso</h1>
      <hr />
      
      <h2>Daftar Employee (Port 8080)</h2>
      <ul>
        {employees.map(emp => (
          <li key={emp.id}>{emp.name} - <b>{emp.position}</b></li>
        ))}
      </ul>

      <h2>Daftar Customer (Port 8081)</h2>
      <ul>
        {customers.map(cust => (
          <li key={cust.id}>{cust.name} ({cust.email})</li>
        ))}
      </ul>
    </div>
  );
}

export default App;