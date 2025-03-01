# Use the official Node.js 20 image as the base image
FROM node:20

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies (including dev dependencies for building)
RUN npm install

# Copy the Prisma schema
COPY prisma ./prisma

# Generate Prisma Client (this doesn't need database connection)
RUN npx prisma generate


# Copy the start script
COPY start.js ./

# Copy the rest of the application
COPY . .

# Use the startup script as the entry point
CMD ["node", "start.js"]
