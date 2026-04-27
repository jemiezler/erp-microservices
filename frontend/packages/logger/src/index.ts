const colors = {
  reset: "\x1b[0m",
  blue: "\x1b[1;34m",
  cyan: "\x1b[0;36m",
  yellow: "\x1b[1;33m",
  green: "\x1b[0;32m",
  red: "\x1b[0;31m",
};

const getTime = () => {
  const now = new Date();
  return now.toLocaleTimeString("en-GB", { hour12: false });
};

export const logger = {
  info: (service: string, message: string) => {
    console.log(`${colors.blue}[${service}] INFO: ${message}${colors.reset}`);
  },
  success: (service: string, message: string) => {
    console.log(`${colors.green}[${service}] SUCCESS: ${message}${colors.reset}`);
  },
  error: (service: string, message: string) => {
    console.error(`${colors.red}[${service}] ERROR: ${message}${colors.reset}`);
  },
  request: (service: string, method: string, path: string, status: number, latency: string) => {
    console.log(`${colors.cyan}[${service}]${colors.reset} ${getTime()} | ${status} | ${latency} | - | ${method} | ${path}`);
  }
};
