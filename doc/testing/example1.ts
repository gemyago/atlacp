// Example TypeScript file: example1.ts
// This file contains a variety of TypeScript constructs for testing purposes.
// Update marker 1

(() => {
// --- Interfaces ---
interface User {
    id: number;
    name: string;
    email: string;
    isActive: boolean;
    roles: string[];
}

interface Product {
    id: number;
    name: string;
    price: number;
    tags?: string[];
}

// --- Enums ---
enum LogLevel {
    DEBUG = 'debug',
    INFO = 'info',
    WARN = 'warn',
    ERROR = 'error',
}

// --- Classes ---
class Logger {
    private level: LogLevel;
    constructor(level: LogLevel) {
        this.level = level;
    }
    log(message: string, level: LogLevel = LogLevel.INFO) {
        if (this.shouldLog(level)) {
            console.log(`[${level}] ${message}`);
        }
    }
    private shouldLog(level: LogLevel): boolean {
        const order = [LogLevel.DEBUG, LogLevel.INFO, LogLevel.WARN, LogLevel.ERROR];
        return order.indexOf(level) >= order.indexOf(this.level);
    }
}
// Update marker 2

// --- Functions ---
function getUserById(users: User[], id: number): User | undefined {
    return users.find(u => u.id === id);
}

function calculateTotal(products: Product[]): number {
    return products.reduce((sum, p) => sum + p.price, 0);
}

function filterActiveUsers(users: User[]): User[] {
    return users.filter(u => u.isActive);
}

function addRoleToUser(user: User, role: string): User {
    return { ...user, roles: [...user.roles, role] };
}

// --- Generics ---
function identity<T>(value: T): T {
    return value;
}

function mapArray<T, U>(arr: T[], fn: (item: T) => U): U[] {
    return arr.map(fn);
}

// --- Utility Types ---
type ReadonlyUser = Readonly<User>;
type PartialProduct = Partial<Product>;

// --- Example Data ---
const users: User[] = [
    { id: 1, name: 'Alice', email: 'alice@example.com', isActive: true, roles: ['admin'] },
    { id: 2, name: 'Bob', email: 'bob@example.com', isActive: false, roles: ['user'] },
    { id: 3, name: 'Charlie', email: 'charlie@example.com', isActive: true, roles: ['user', 'editor'] },
];
// Update marker 3

const products: Product[] = [
    { id: 1, name: 'Laptop', price: 1200, tags: ['electronics', 'computers'] },
    { id: 2, name: 'Phone', price: 800 },
    { id: 3, name: 'Book', price: 20, tags: ['education'] },
];

// --- Main Logic ---
const logger = new Logger(LogLevel.DEBUG);

logger.log('Application started', LogLevel.INFO);

const activeUsers = filterActiveUsers(users);
logger.log(`Active users: ${activeUsers.map(u => u.name).join(', ')}`);

const total = calculateTotal(products);
logger.log(`Total price: $${total}`);

const alice = getUserById(users, 1);
if (alice) {
    logger.log(`Found user: ${alice.name}`);
}

const newUser = addRoleToUser(users[1], 'editor');
logger.log(`Updated user roles: ${newUser.roles.join(', ')}`);

const idVal = identity<string>('test');
logger.log(`Identity function returned: ${idVal}`);

const mapped = mapArray<number, string>([1, 2, 3], n => `Number: ${n}`);
logger.log(`Mapped array: ${mapped.join('; ')}`);

// --- File Operations (Mocked) ---
function saveDataToFile(filename: string, data: string) {
    // Simulate file save
    logger.log(`Data would be saved to ${filename}`);
}

saveDataToFile('output.txt', JSON.stringify(users));
// Update marker 4

// --- More Example Code ---
for (let i = 0; i < 10; i++) {
    logger.log(`Loop iteration: ${i}`);
}

function* idGenerator() {
    let id = 0;
    while (true) {
        yield id++;
    }
}

const gen = idGenerator();
for (let i = 0; i < 5; i++) {
    logger.log(`Generated id: ${gen.next().value}`);
}

class Calculator {
    add(a: number, b: number): number {
        return a + b;
    }
    multiply(a: number, b: number): number {
        return a * b;
    }
}

const calc = new Calculator();
logger.log(`Calc add: ${calc.add(2, 3)}`);
logger.log(`Calc multiply: ${calc.multiply(4, 5)}`);

// --- End of File ---
// Update marker 5
})();