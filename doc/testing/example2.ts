// Example TypeScript file: example2.ts
// This file demonstrates more advanced TypeScript features and patterns.
// Update marker 1

// --- Types and Interfaces ---
type UUID = string;

interface Task {
    id: UUID;
    title: string;
    completed: boolean;
    dueDate?: Date;
}

interface Project {
    id: UUID;
    name: string;
    tasks: Task[];
}

// --- Enums ---
enum TaskStatus {
    TODO = 'todo',
    IN_PROGRESS = 'in_progress',
    DONE = 'done',
}

// --- Classes ---
class TaskManager {
    private tasks: Map<UUID, Task> = new Map();
    private listeners: { [event: string]: Function[] } = {};
    // Update marker 2
    addTask(task: Task) {
        this.tasks.set(task.id, task);
        this.emit('taskAdded', task);
    }
    completeTask(id: UUID) {
        const task = this.tasks.get(id);
        if (task) {
            task.completed = true;
            this.emit('taskCompleted', task);
        }
    }
    getTasks(): Task[] {
        return Array.from(this.tasks.values());
    }
    getPendingTasks(): Task[] {
        return this.getTasks().filter(t => !t.completed);
    }
    removeTask(id: UUID) {
        this.tasks.delete(id);
        this.emit('taskRemoved', id);
    }
    on(event: string, listener: Function) {
        if (!this.listeners[event]) {
            this.listeners[event] = [];
        }
        this.listeners[event].push(listener);
    }
    emit(event: string, ...args: any[]) {
        if (this.listeners[event]) {
            for (const listener of this.listeners[event]) {
                listener(...args);
            }
        }
    }
}

// --- Utility Functions ---
function generateUUID(): UUID {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, c => {
        const r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

function createTask(title: string, dueDate?: Date): Task {
    return {
        id: generateUUID(),
        title,
        completed: false,
        dueDate,
    };
}

function createProject(name: string): Project {
    return {
        id: generateUUID(),
        name,
        tasks: [],
    };
}

// --- Example Data ---
const project = createProject('AI Project');
const task1 = createTask('Design API', new Date('2024-07-01'));
const task2 = createTask('Implement backend');
const task3 = createTask('Write tests', new Date('2024-07-10'));
project.tasks.push(task1, task2, task3);

// --- Task Manager Usage ---
const manager = new TaskManager();
manager.on('taskAdded', (task: Task) => {
    console.log(`Task added: ${task.title}`);
});
manager.on('taskCompleted', (task: Task) => {
    console.log(`Task completed: ${task.title}`);
});
manager.on('taskRemoved', (id: UUID) => {
    console.log(`Task removed: ${id}`);
});

for (const t of project.tasks) {
    manager.addTask(t);
}

manager.completeTask(task1.id);
manager.removeTask(task2.id);

console.log('Pending tasks:', manager.getPendingTasks().map(t => t.title));

// --- Advanced Types ---
type TaskMap = Record<UUID, Task>;
type ProjectWithStatus = Project & { status: TaskStatus };

// --- Example of Type Guards ---
function isProject(obj: any): obj is Project {
    return obj && typeof obj.id === 'string' && Array.isArray(obj.tasks);
}

// --- Example of Promises and Async/Await ---
async function fetchProject(id: UUID): Promise<Project> {
    // Simulate async fetch
    // Update marker 3
    return new Promise(resolve => {
        setTimeout(() => {
            resolve(createProject('Fetched Project'));
        }, 100);
    });
}

// --- Example of Generics ---
function wrapInArray<T>(item: T): T[] {
    return [item];
}

// --- Example of Decorators ---
function logClass(target: Function) {
    console.log(`Class created: ${target.name}`);
}

@logClass
class Service {
    name: string;
    constructor(name: string) {
        this.name = name;
    }
    start() {
        console.log(`${this.name} started`);
    }
}

const service = new Service('NotificationService');
service.start();

// --- Example of Namespaces ---
namespace MathUtils {
    export function add(a: number, b: number): number {
        return a + b;
    }
    export function multiply(a: number, b: number): number {
        return a * b;
    }
}

console.log('Math add:', MathUtils.add(2, 3));
console.log('Math multiply:', MathUtils.multiply(4, 5));

// --- End of File ---
// Update marker 4
