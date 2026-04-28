'use client';

import { useState, useEffect } from 'react';
import { employeeApi, leaveApi, attendanceApi, payrollApi, Employee, Leave, Attendance, AttendanceStats, Payroll, ApiError } from '@/lib/api-client';
import { 
  Button, 
  Card, 
  CardContent, 
  CardHeader, 
  CardTitle, 
  Table, 
  TableBody, 
  TableCell, 
  TableHead, 
  TableHeader, 
  TableRow,
  Badge,
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger
} from '@erp/ui';

interface DashboardStats {
  totalEmployees: number;
  activeEmployees: number;
  pendingLeaves: number;
  todayPresent: number;
  pendingPayrolls: number;
}

export default function HRDashboard() {
  const [stats, setStats] = useState<DashboardStats>({
    totalEmployees: 0,
    activeEmployees: 0,
    pendingLeaves: 0,
    todayPresent: 0,
    pendingPayrolls: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'overview' | 'employees' | 'leaves' | 'attendance' | 'payroll'>('overview');

  useEffect(() => {
    const loadDashboard = async () => {
      try {
        setLoading(true);
        setError(null);

        const [empRes, pendingLeavesRes, statsRes, payrollRes] = await Promise.all([
          employeeApi.getAll(1, 1),
          leaveApi.getPending(),
          attendanceApi.getStats(),
          payrollApi.getPending(),
        ]);

        setStats({
          totalEmployees: empRes.pagination?.total || 0,
          activeEmployees: empRes.data?.filter((e: any) => e.status === 'active').length || 0,
          pendingLeaves: pendingLeavesRes.pagination?.total || 0,
          todayPresent: (statsRes.data as any)?.present || 0,
          pendingPayrolls: payrollRes.pagination?.total || 0,
        });
      } catch (err) {
        const apiError = err as ApiError;
        setError(apiError.message || 'Failed to load dashboard');
        console.error('Dashboard error:', err);
      } finally {
        setLoading(false);
      }
    };

    loadDashboard();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-slate-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-slate-600">Loading dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <header className="bg-white border-b border-slate-200 shadow-sm">
        <div className="max-w-7xl mx-auto px-6 py-6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-slate-900">HR Dashboard</h1>
              <p className="text-slate-600 mt-1">Manage employees, leaves, attendance & payroll</p>
            </div>
            <div className="bg-blue-100 text-blue-800 px-4 py-2 rounded-full text-sm font-semibold">
              HR-MFE v1.0
            </div>
          </div>
        </div>
      </header>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mx-6 mt-6">
          <p className="font-semibold">Error: {error}</p>
        </div>
      )}

      <div className="max-w-7xl mx-auto px-6 py-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-8">
          <Card className="border-l-4 border-blue-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-slate-600 text-sm font-medium uppercase">Total Employees</CardTitle>
            </CardHeader>
            <CardContent className="flex items-center justify-between">
              <p className="text-2xl font-bold text-slate-900">{stats.totalEmployees}</p>
              <div className="text-3xl text-blue-200">👥</div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-green-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-slate-600 text-sm font-medium uppercase">Active</CardTitle>
            </CardHeader>
            <CardContent className="flex items-center justify-between">
              <p className="text-2xl font-bold text-slate-900">{stats.activeEmployees}</p>
              <div className="text-3xl text-green-200">✓</div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-yellow-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-slate-600 text-sm font-medium uppercase">Pending Leaves</CardTitle>
            </CardHeader>
            <CardContent className="flex items-center justify-between">
              <p className="text-2xl font-bold text-slate-900">{stats.pendingLeaves}</p>
              <div className="text-3xl text-yellow-200">⏳</div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-purple-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-slate-600 text-sm font-medium uppercase">Present Today</CardTitle>
            </CardHeader>
            <CardContent className="flex items-center justify-between">
              <p className="text-2xl font-bold text-slate-900">{stats.todayPresent}</p>
              <div className="text-3xl text-purple-200">📍</div>
            </CardContent>
          </Card>

          <Card className="border-l-4 border-orange-500">
            <CardHeader className="pb-2">
              <CardTitle className="text-slate-600 text-sm font-medium uppercase">Pending Payroll</CardTitle>
            </CardHeader>
            <CardContent className="flex items-center justify-between">
              <p className="text-2xl font-bold text-slate-900">{stats.pendingPayrolls}</p>
              <div className="text-3xl text-orange-200">💰</div>
            </CardContent>
          </Card>
        </div>

        <Tabs defaultValue="overview" className="w-full">
          <TabsList className="grid w-full grid-cols-5">
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="employees">Employees</TabsTrigger>
            <TabsTrigger value="leaves">Leaves</TabsTrigger>
            <TabsTrigger value="attendance">Attendance</TabsTrigger>
            <TabsTrigger value="payroll">Payroll</TabsTrigger>
          </TabsList>
          
          <TabsContent value="overview" className="mt-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card className="bg-blue-50/50 border-blue-200">
                <CardHeader>
                  <CardTitle className="text-blue-900 flex items-center gap-2">
                    🎯 Quick Actions
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-2 gap-4">
                    <Button variant="outline" className="justify-start">Add Employee</Button>
                    <Button variant="outline" className="justify-start">Review Leaves</Button>
                    <Button variant="outline" className="justify-start">Run Payroll</Button>
                    <Button variant="outline" className="justify-start">Export Data</Button>
                  </div>
                </CardContent>
              </Card>
              <Card className="bg-green-50/50 border-green-200">
                <CardHeader>
                  <CardTitle className="text-green-900 flex items-center gap-2">
                    📊 System Status
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <ul className="space-y-2 text-sm text-green-800">
                    <li className="flex items-center gap-2">✓ <span className="font-medium">Database:</span> Connected</li>
                    <li className="flex items-center gap-2">✓ <span className="font-medium">HR Service:</span> Running</li>
                    <li className="flex items-center gap-2">✓ <span className="font-medium">All modules:</span> Active</li>
                  </ul>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          <TabsContent value="employees">
            <Card>
              <CardContent className="pt-6">
                <EmployeesTab />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="leaves">
            <Card>
              <CardContent className="pt-6">
                <LeavesTab />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="attendance">
            <Card>
              <CardContent className="pt-6">
                <AttendanceTab />
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="payroll">
            <Card>
              <CardContent className="pt-6">
                <PayrollTab />
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}

function EmployeesTab() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadEmployees = async () => {
      try {
        const res = await employeeApi.getAll();
        setEmployees(res.data || []);
      } catch (err) {
        console.error('Failed to load employees:', err);
      } finally {
        setLoading(false);
      }
    };
    loadEmployees();
  }, []);

  if (loading) return <div className="text-center py-8">Loading employees...</div>;

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Employee ID</TableHead>
          <TableHead>Name</TableHead>
          <TableHead>Email</TableHead>
          <TableHead>Position</TableHead>
          <TableHead>Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {employees.map((emp: any) => (
          <TableRow key={emp.id}>
            <TableCell className="font-mono text-xs">{emp.employee_id}</TableCell>
            <TableCell className="font-medium">{emp.name}</TableCell>
            <TableCell>{emp.email}</TableCell>
            <TableCell>{emp.position}</TableCell>
            <TableCell>
              <Badge variant={emp.status === 'active' ? 'default' : 'secondary'}>
                {emp.status}
              </Badge>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function LeavesTab() {
  const [leaves, setLeaves] = useState<Leave[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadLeaves = async () => {
      try {
        const res = await leaveApi.getPending();
        setLeaves(res.data || []);
      } catch (err) {
        console.error('Failed to load leaves:', err);
      } finally {
        setLoading(false);
      }
    };
    loadLeaves();
  }, []);

  if (loading) return <div className="text-center py-8">Loading leave requests...</div>;

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Employee</TableHead>
          <TableHead>Leave Type</TableHead>
          <TableHead>Start Date</TableHead>
          <TableHead>End Date</TableHead>
          <TableHead>Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {leaves.map((leave: any) => (
          <TableRow key={leave.id}>
            <TableCell className="font-medium">{leave.employee_id}</TableCell>
            <TableCell>{leave.leave_type}</TableCell>
            <TableCell>{new Date(leave.start_date).toLocaleDateString()}</TableCell>
            <TableCell>{new Date(leave.end_date).toLocaleDateString()}</TableCell>
            <TableCell>
              <Badge variant={
                leave.status === 'approved' ? 'default' : 
                leave.status === 'pending' ? 'outline' : 'destructive'
              }>
                {leave.status}
              </Badge>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function AttendanceTab() {
  const [stats, setStats] = useState<AttendanceStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadStats = async () => {
      try {
        const res = await attendanceApi.getStats();
        setStats(res.data as any);
      } catch (err) {
        console.error('Failed to load attendance stats:', err);
      } finally {
        setLoading(false);
      }
    };
    loadStats();
  }, []);

  if (loading) return <div className="text-center py-8">Loading attendance data...</div>;

  if (!stats) return <div className="text-center py-8 text-red-600">Failed to load attendance</div>;

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card className="bg-green-50/30 border-green-100">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-slate-500">Present</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-bold text-green-600">{stats.present}</p>
        </CardContent>
      </Card>
      <Card className="bg-red-50/30 border-red-100">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-slate-500">Absent</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-bold text-red-600">{stats.absent}</p>
        </CardContent>
      </Card>
      <Card className="bg-yellow-50/30 border-yellow-100">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-slate-500">On Leave</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-bold text-yellow-600">{stats.on_leave}</p>
        </CardContent>
      </Card>
      <Card className="bg-blue-50/30 border-blue-100">
        <CardHeader className="pb-2">
          <CardTitle className="text-sm font-medium text-slate-500">Attendance %</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-bold text-blue-600">{stats.attendance_percentage.toFixed(1)}%</p>
        </CardContent>
      </Card>
    </div>
  );
}

function PayrollTab() {
  const [payrolls, setPayrolls] = useState<Payroll[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadPayrolls = async () => {
      try {
        const res = await payrollApi.getPending();
        setPayrolls(res.data || []);
      } catch (err) {
        console.error('Failed to load payrolls:', err);
      } finally {
        setLoading(false);
      }
    };
    loadPayrolls();
  }, []);

  if (loading) return <div className="text-center py-8">Loading payroll data...</div>;

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Employee</TableHead>
          <TableHead>Month</TableHead>
          <TableHead className="text-right">Gross</TableHead>
          <TableHead className="text-right">Deductions</TableHead>
          <TableHead className="text-right">Net</TableHead>
          <TableHead>Status</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {payrolls.map((pr: any) => (
          <TableRow key={pr.id}>
            <TableCell className="font-medium">{pr.employee_id}</TableCell>
            <TableCell>{pr.month}</TableCell>
            <TableCell className="text-right">${pr.gross_salary.toLocaleString()}</TableCell>
            <TableCell className="text-right text-red-600">-${pr.deductions.toLocaleString()}</TableCell>
            <TableCell className="text-right font-bold text-green-600">${pr.net_salary.toLocaleString()}</TableCell>
            <TableCell>
              <Badge variant={pr.status === 'approved' ? 'default' : 'outline'}>
                {pr.status}
              </Badge>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
