"use client";

import { useEffect, useState } from 'react';

interface Employee {
  ID: number;
  name: string;
  email: string;
  position: string;
  status: string;
}

export default function HRDashboard() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('http://localhost:8080/api/v1/hr/employees', {
      headers: {
        'Authorization': 'Bearer dev-token'
      }
    })
      .then(res => res.json())
      .then(data => {
        setEmployees(data || []);
        setLoading(false);
      })
      .catch(err => {
        console.error('Failed to fetch employees:', err);
        setLoading(false);
      });
  }, []);

  return (
    <div className="p-8 bg-slate-50 min-h-screen font-sans">
      <div className="max-w-6xl mx-auto">
        <header className="mb-8 border-b pb-4 flex justify-between items-center">
          <h1 className="text-2xl font-bold text-slate-800">HR Micro-Frontend Dashboard</h1>
          <span className="bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded">SERVICE: HR-MFE</span>
        </header>

        {loading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-6">
            <div className="bg-white shadow rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-slate-200">
                <thead className="bg-slate-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Name</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Email</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Position</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-slate-500 uppercase tracking-wider">Status</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-slate-200 text-sm">
                  {employees.map((emp) => (
                    <tr key={emp.ID}>
                      <td className="px-6 py-4 whitespace-nowrap font-medium text-slate-900">{emp.name}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-slate-600">{emp.email}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-slate-600">{emp.position}</td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${emp.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}`}>
                          {emp.status}
                        </span>
                      </td>
                    </tr>
                  ))}
                  {employees.length === 0 && (
                    <tr>
                      <td colSpan={4} className="px-6 py-12 text-center text-slate-500 italic">No employees found in database.</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
