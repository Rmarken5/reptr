# Use the official MongoDB image as the base image
FROM mongo:7.0-jammy

# Expose the default MongoDB port (27017)
EXPOSE 27017

# Set a custom data directory (optional)
# ENV MONGO_DATA_DIR /data/db

# Set a custom log directory (optional)
# ENV MONGO_LOG_DIR /var/log/mongodb

# Optionally, you can set environment variables for MongoDB configuration
# For example, if you want to enable authentication, you can set the root user and password
# ENV MONGO_INITDB_ROOT_USERNAME myrootuser
# ENV MONGO_INITDB_ROOT_PASSWORD myrootpassword

# Entry point for the container (you don't need to specify this if you are using the official image)
# CMD ["mongod"]

# Optionally, if you have a custom MongoDB configuration file (e.g., mongod.conf) you want to use, you can copy it:
COPY ./mongo/mongod.conf /etc/mongod.conf
