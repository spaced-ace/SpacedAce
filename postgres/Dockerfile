FROM postgres:16.2-bookworm

# Install necessary packages and dependencies
RUN apt-get update && apt-get install -y \
    curl \
    gcc \
    make \
    pkg-config \
    autoconf \
    libssl-dev \
    libpq5 \
    libc6 \
    git \
    postgresql-server-dev-16 \
    && rm -rf /var/lib/apt/lists/*

# Install Rust (required for pg_cron installation)
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

# Clone and build pgvector from source
RUN git clone https://github.com/pgvector/pgvector.git /pgvector
WORKDIR /pgvector
RUN make && make install

# Install pg_cron using pg-trunk
RUN cargo install pg-trunk
RUN trunk install pg_cron

# Configure PostgreSQL to preload pg_cron
RUN echo "shared_preload_libraries = 'pg_cron'" >> /usr/share/postgresql/postgresql.conf.sample && cp /usr/share/postgresql/postgresql.conf.sample /etc/postgresql/postgresql.conf

# Switch back to the default work directory
WORKDIR /
