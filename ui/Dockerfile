FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Expose port 3002 (to match the dev server port)
EXPOSE 3002

# Set environment variable to production
ENV NODE_ENV=production

# Start the development server (which includes proxy support)
CMD ["npm", "start"] 