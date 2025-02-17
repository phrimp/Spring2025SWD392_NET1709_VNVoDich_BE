const { execSync } = require('child_process');
const { setTimeout } = require('timers/promises');

async function waitForDatabase() {
    console.log("Waiting for database to be ready...");
    let retryCount = 0;
    const maxRetries = 30;

    while (retryCount < maxRetries) {
        try {
            // Force reset and push the new schema
            console.log("Pushing new schema to database (with reset)...");
            execSync('npx prisma db push --force-reset --accept-data-loss', { stdio: 'inherit' });
            
            // Generate the Prisma client
            console.log("Generating Prisma Client...");
            execSync('npx prisma generate', { stdio: 'inherit' });
            
            console.log("Database schema is ready!");
            return true;
        } catch (error) {
            retryCount++;
            console.log(`\nAttempt ${retryCount}/${maxRetries} failed`);
            console.log("Error:", error.message);
            
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
