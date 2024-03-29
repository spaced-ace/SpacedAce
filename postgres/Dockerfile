FROM postgres:16.2-bookworm

RUN apt-get update && apt-get install -y \
	curl \
	gcc \
	make \
	pkg-config \
	autoconf \
	libssl-dev \
	libpq5 \
	libc6 \
	&& rm -rf /var/lib/apt/lists/*
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
RUN ./root/.cargo/bin/cargo install pg-trunk
RUN ./root/.cargo/bin/trunk install pgvector
RUN ./root/.cargo/bin/trunk install pg_cron
RUN echo "shared_preload_libraries = 'pg_cron'" >> /usr/share/postgresql/postgresql.conf.sample && cp /usr/share/postgresql/postgresql.conf.sample /etc/postgresql/postgresql.conf


