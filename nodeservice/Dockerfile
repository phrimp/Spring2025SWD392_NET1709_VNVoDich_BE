# Use the official Node.js 20 image as the base image
FROM node:20

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Copy .env file (optional if required during build)
COPY .env .env

# Install dependencies (including dev dependencies for building)
RUN npm install

# Copy the Prisma schema and generate the Prisma Client
COPY prisma ./prisma
RUN npx prisma generate
RUN npx prisma db push

# Copy the rest of the application
COPY . .

# Start the application
CMD ["npm", "run", "dev"]