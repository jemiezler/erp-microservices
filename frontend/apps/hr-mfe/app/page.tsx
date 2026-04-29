'use client';

import { useState, useEffect } from 'react';
import { employeeApi, leaveApi, attendanceApi, payrollApi, Employee, AttendanceStats } from '../lib/api-client';
import { 
  Button, 
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
  TabsTrigger,
  BentoGrid,
  BentoCard,
  cn
} from '@erp/ui';
import { Users, Calendar, Clock, CreditCard, ChevronRight, Activity, Plus } from 'lucide-react';

export function HRDashboardContent() {
  const [stats, setStats] = useState({
    totalEmployees: 0,
    activeEmployees: 0,
    pendingLeaves: 0,
    todayPresent: 0,
    pendingPayrolls: 0,
    attendanceRate: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadDashboard = async () => {
      try {
        setLoading(true);
        const [empRes, pendingLeavesRes, statsRes, payrollRes] = await Promise.all([
          employeeApi.getAll(1, 1),
          leaveApi.getPending(),
          attendanceApi.getStats(),
          payrollApi.getPending(),
        ]);

        const attendanceData = statsRes.data as AttendanceStats;
        setStats({
          totalEmployees: empRes.pagination?.total || 0,
          activeEmployees: empRes.data?.filter((e: Employee) => e.status === 'active').length || 0,
          pendingLeaves: pendingLeavesRes.pagination?.total || 0,
          todayPresent: attendanceData?.present || 0,
          pendingPayrolls: payrollRes.pagination?.total || 0,
          attendanceRate: attendanceData?.attendance_percentage || 0,
        });
      } catch (err) {
        console.error('Dashboard error:', err);
      } finally {
        setLoading(false);
      }
    };
    void loadDashboard();
  }, []);

  if (loading) {
    return (
      <div className='flex items-center justify-center py-40'>
        <div className='w-12 h-12 border-4 border-primary border-t-transparent rounded-full animate-spin' />
      </div>
    );
  }

  return (
    <BentoGrid className='p-0 gap-6'>
      {/* Metrics Row */}
      <BentoCard span={3} className='bg-card text-card-foreground border-border'>
           <div className='flex items-center gap-3 mb-4'>
              <div className='p-2 bg-secondary rounded-lg'><Users size={18} className='text-primary' /></div>
              <span className='text-[10px] font-bold uppercase tracking-widest text-muted-foreground'>Total Staff</span>
           </div>
           <div className='text-4xl font-black text-foreground'>{stats.totalEmployees}</div>
           <div className='mt-4 flex items-center gap-1 text-[10px] font-bold text-success-foreground'>
              <Activity size={12} /> +2.4% vs last month
           </div>
      </BentoCard>

      <BentoCard span={3} className='bg-card text-card-foreground border-border'>
           <div className='flex items-center gap-3 mb-4'>
              <div className='p-2 bg-success rounded-lg'><Clock size={18} className='text-success-foreground' /></div>
              <span className='text-[10px] font-bold uppercase tracking-widest text-muted-foreground'>Attendance</span>
           </div>
           <div className='text-4xl font-black text-foreground'>{stats.attendanceRate.toFixed(1)}%</div>
           <div className='mt-4 text-[10px] font-bold text-muted-foreground'>{stats.todayPresent} present today</div>
      </BentoCard>

      <BentoCard span={3} className='bg-card text-card-foreground border-border'>
           <div className='flex items-center gap-3 mb-4'>
              <div className='p-2 bg-alert rounded-lg'><Calendar size={18} className='text-alert-foreground' /></div>
              <span className='text-[10px] font-bold uppercase tracking-widest text-muted-foreground'>Leaves</span>
           </div>
           <div className='text-4xl font-black text-foreground'>{stats.pendingLeaves}</div>
           <div className='mt-4 text-[10px] font-bold text-alert-foreground'>Requires Approval</div>
      </BentoCard>

      <BentoCard span={3} className='bg-primary text-primary-foreground border-border'>
           <div className='flex items-center gap-3 mb-4'>
              <div className='p-2 bg-primary-foreground/10 rounded-lg'><CreditCard size={18} className='text-primary-foreground' /></div>
              <span className='text-[10px] font-bold uppercase tracking-widest text-primary-foreground/40'>Payroll</span>
           </div>
           <div className='text-4xl font-black'>{stats.pendingPayrolls}</div>
           <div className='mt-4 text-[10px] font-bold text-success'>Ready to process</div>
      </BentoCard>

      {/* Interface Breakout - Main Content */}
      <BentoCard span={12} className='bg-card text-card-foreground border-border p-0 overflow-hidden'>
        <Tabs defaultValue='employees' className='w-full'>
          <div className='px-8 pt-8 flex items-center justify-between border-b border-border pb-4'>
            <TabsList className='bg-muted p-1 rounded-full'>
              <TabsTrigger value='employees' className='rounded-full px-6 data-[state=active]:bg-card data-[state=active]:text-card-foreground data-[state=active]:shadow-sm'>Staff</TabsTrigger>
              <TabsTrigger value='leaves' className='rounded-full px-6 data-[state=active]:bg-card data-[state=active]:text-card-foreground data-[state=active]:shadow-sm'>Requests</TabsTrigger>
              <TabsTrigger value='payroll' className='rounded-full px-6 data-[state=active]:bg-card data-[state=active]:text-card-foreground data-[state=active]:shadow-sm'>Finance</TabsTrigger>
            </TabsList>
            <Button size='sm' className='bg-primary text-primary-foreground rounded-full flex items-center gap-2'>
              <Plus size={16} /> New Record
            </Button>
          </div>
          
          <TabsContent value='employees' className='p-8'>
             <EmployeesTab />
          </TabsContent>
          <TabsContent value='leaves' className='p-8'>
             <LeavesTab />
          </TabsContent>
          <TabsContent value='payroll' className='p-8'>
             <PayrollTab />
          </TabsContent>
        </Tabs>
      </BentoCard>
    </BentoGrid>
  );
}

function EmployeesTab() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  useEffect(() => {
    const loadEmployees = async () => {
      try {
        const res = await employeeApi.getAll();
        setEmployees(res.data || []);
      } catch (err) {
        console.error('Failed to load employees:', err);
      }
    };

    void loadEmployees();
  }, []);

  return (
    <div className='rounded-[24px] border border-border overflow-hidden bg-card shadow-sm'>
      <Table>
        <TableHeader className='bg-muted/50'>
          <TableRow className='border-0'>
            <TableHead className='font-bold text-muted-foreground text-[10px] uppercase'>Identity</TableHead>
            <TableHead className='font-bold text-muted-foreground text-[10px] uppercase'>Role</TableHead>
            <TableHead className='font-bold text-muted-foreground text-[10px] uppercase'>Status</TableHead>
            <TableHead className='text-right'></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {employees.map((emp: Employee) => (
            <TableRow key={emp.id} className='hover:bg-muted/30 transition-colors border-border'>
              <TableCell>
                <div className='flex items-center gap-3'>
                   <div className='w-8 h-8 rounded-full bg-secondary flex items-center justify-center font-bold text-[10px] text-secondary-foreground'>
                      {emp.name.charAt(0)}
                   </div>
                   <div>
                      <div className='font-bold text-sm text-foreground'>{emp.name}</div>
                      <div className='text-[10px] text-muted-foreground font-medium'>{emp.email}</div>
                   </div>
                </div>
              </TableCell>
              <TableCell className='text-xs font-medium text-foreground'>{emp.position}</TableCell>
              <TableCell>
                <Badge className={cn("rounded-full border-0 px-3", 
                  emp.status.toLowerCase() === 'active' ? 'bg-success text-success-foreground' : 
                  emp.status.toLowerCase() === 'on leave' ? 'bg-alert text-alert-foreground' :
                  'bg-secondary text-secondary-foreground'
                )}>
                  {emp.status}
                </Badge>
              </TableCell>
              <TableCell className='text-right'>
                <Button variant='ghost' size='icon' className='rounded-full text-muted-foreground hover:text-foreground'><ChevronRight size={16} /></Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}

// Simplified placeholders for other tabs to keep focus on grid structure
function LeavesTab() { return <div className='p-20 text-center text-muted-foreground font-bold'>Leave Requests Interface</div> }
function PayrollTab() { return <div className='p-20 text-center text-muted-foreground font-bold'>Payroll Engine Interface</div> }

export default function HRDashboardPage() {
  return (
    <div className='min-h-screen bg-background font-sans text-foreground'>
      <header className='px-10 py-10 flex flex-col md:flex-row md:items-end justify-between gap-6'>
        <div className='flex flex-col items-start'>
           <Badge className='bg-primary text-primary-foreground rounded-full mb-4 px-3 py-1'>Core Module</Badge>
           <h1 className='text-6xl font-black tracking-tighter text-foreground'>Human Resources</h1>
        </div>
        <div className='flex items-center gap-4 bg-card/50 p-2 rounded-3xl border border-border backdrop-blur-sm shrink-0'>
           <div className='text-right px-4'>
              <div className='text-[10px] font-bold text-muted-foreground uppercase tracking-widest'>System Health</div>
              <div className='text-sm font-bold text-success-foreground'>Stable 100%</div>
           </div>
           <div className='w-12 h-12 rounded-2xl bg-primary flex items-center justify-center shadow-lg'>
              <Activity className='text-primary-foreground' size={24} />
           </div>
        </div>
      </header>
      <main className='px-10 pb-20'>
        <HRDashboardContent />
      </main>
    </div>
  );
}
