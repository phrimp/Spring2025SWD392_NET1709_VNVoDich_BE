const { execSync } = require('child_process');
const { setTimeout } = require('timers/promises');

async function waitForDatabase() {
    console.log("Waiting for database to be ready...");
    let retryCount = 0;
    const maxRetries = 30;

    while (retryCount < maxRetries) {
        try {
            // First pull the existing schema
            console.log("Pulling current database schema...");
            execSync('npx prisma db pull', { stdio: 'inherit' });
            
            // Then push any new changes without forcing data loss
            console.log("Safely pushing schema updates...");
            execSync('npx prisma generate', { stdio: 'inherit' });
            
            console.log("Database schema is ready!");
            return true;
        } catch (error) {
            retryCount++;
            console.log(`\nAttempt ${retryCount}/${maxRetries} failed`);
            
            if (retryCount === maxRetries) {
                console.error(`Failed to sync schema after ${maxRetries} attempts`);
                process.exit(1);
            }
            console.log(`Waiting 5 seconds before next attempt...\n`);
            await setTimeout(5000);
        }
    }
}

async function start() {
    await waitForDatabase();
    console.log("Starting application...");
    execSync('npm run dev', { stdio: 'inherit' });
}

start().catch(error => {
    console.error('Startup failed:', error);
    process.exit(1);
});
