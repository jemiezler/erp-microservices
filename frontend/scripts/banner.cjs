const pkg = JSON.parse(require('fs').readFileSync(process.argv[2], 'utf8'));
const port = process.argv[3] || 3000;
const name = pkg.name;
const version = pkg.dependencies.next;

const width = 51;
const line = '─'.repeat(width);

console.log('┌' + line + '┐');
console.log('│' + 'Next.js ' + version.padStart(width - 8) + ' │');
console.log('│' + ('http://127.0.0.1:' + port).padStart(Math.floor((width + ('http://127.0.0.1:' + port).length) / 2)).padEnd(width) + ' │');
console.log('│' + ('(bound on host 0.0.0.0 and port ' + port + ')').padStart(Math.floor((width + ('(bound on host 0.0.0.0 and port ' + port + ')').length) / 2)).padEnd(width) + ' │');
console.log('│' + ' '.repeat(width) + ' │');
console.log('│ App ................ ' + name.padEnd(width - 23) + ' │');
console.log('│ PID ................ ' + process.pid.toString().padEnd(width - 23) + ' │');
console.log('└' + line + '┘');
