FROM debian:jessie
MAINTAINER gubin <gubin@empower.cn>

#苏州远鼎提供的代码
ENV ODOO_VERSION yd-10.0

# #更换源
# RUN echo "deb http://mirrors.ustc.edu.cn/debian jessie main contrib non-free\n \
#     deb-src http://mirrors.ustc.edu.cn/debian jessie main contrib non-free\n \
#     deb http://mirrors.ustc.edu.cn/debian jessie-proposed-updates main contrib non-free\n \
#     deb-src http://mirrors.ustc.edu.cn/debian jessie-proposed-updates main contrib non-free\n \
#     deb http://mirrors.ustc.edu.cn/debian jessie-updates main contrib non-free\n \
#     deb-src http://mirrors.ustc.edu.cn/debian jessie-updates main contrib non-free" \
# | tee /etc/apt/sources.list

# Install some deps, lessc and less-plugin-clean-css, and wkhtmltopdf
RUN set -x; \
        apt-get update \
        && apt-get install -y --no-install-recommends \
            ca-certificates \
            curl \
            python-pyinotify \
            python-renderpm \
            python-support \
						libxml2-dev \
						libxslt1-dev \
						libldap2-dev \
						libjpeg-dev \
						python-dev \
						npm \
						build-essential \
						fontconfig \
						libfontconfig1 \
						libx11-6 \
						libxext6 \
						libxrender1 \
						xfonts-base \
						xfonts-75dpi \
						libssl-dev \
						libjpeg62-turbo \
						ttf-wqy-zenhei \
						ttf-wqy-microhei \
						git \
            python \
						python-pip \
						libsasl2-dev \
						libssl-dev \
						libpq-dev \
        && npm install npm -g \
				&& npm cache clean -f \
				&& npm install -g n \
				&& n lts

COPY ./pgdg.list /etc/apt/sources.list.d

RUN curl -s -o - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
		&& apt-get update \
		&& npm install -g less less-plugin-clean-css \
		&& apt-get install -y postgresql-client-10 \
		&& curl -o  wkhtmltox.deb -SL http://nightly.odoo.com/extra/wkhtmltox-0.12.1.2_linux-jessie-amd64.deb \
		&& echo '40e8b906de658a2221b15e4e8cd82565a47d7ee8 wkhtmltox.deb' | sha1sum -c - \
		&& dpkg --force-depends -i wkhtmltox.deb \
		&& apt-get -y install -f --no-install-recommends \
		&& apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false -o APT::AutoRemove::SuggestsImportant=false npm \
    && rm -rf /var/lib/apt/lists/* wkhtmltox.deb


WORKDIR /usr/src

# create the odoo user
RUN adduser --home=/var/lib/odoo --disabled-password --gecos "" --shell=/bin/bash odoo \
		&& pip install --upgrade pip -i https://pypi.douban.com/simple

COPY ./odoo /usr/src
COPY ./sa_addons /usr/src

RUN pip install --upgrade setuptools \
    && rm -rf /usr/src/odoo/.git \
    && rm -rf /usr/src/sa_addons/.git \
    && rm -rf /usr/src/sa_addons/docker \
    && rm -rf /usr/src/sa_addons/.circleci \
    && pip install -r /usr/src/odoo/requirements.txt -i https://pypi.douban.com/simple \
    && pip install -r /usr/src/sa_addons/requirements.txt -i https://pypi.douban.com/simple \
		&& ln -s /usr/src/odoo/setup/odoo /usr/local/bin/odoo \
    && chmod +x /usr/local/bin/odoo \
		&& sed -i "s/fonts\.googleapis\.com/fonts.lug.ustc.edu.cn/g" `grep 'fonts\.googleapis\.com' -rl /usr/src/odoo/odoo/addons`

# Copy entrypoint script and Odoo configuration file
COPY ./entrypoint.sh /
COPY ./odoo.conf /etc/odoo/

# Mount /var/lib/odoo to allow restoring filestore and /mnt/extra-addons for users addons
RUN mkdir -p /mnt/extra-addons \
    && chown -R odoo /mnt/extra-addons \
    && chown odoo /etc/odoo/odoo.conf \
    && chown odoo /entrypoint.sh \
    && chown -R odoo /usr/src

### Remove all dev depends
RUN set -x; \
				apt-get -y install -f --no-install-recommends \
				&& apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false -o APT::AutoRemove::SuggestsImportant=false git \
						build-essential \
						git \
						# python-pip \
						# libsasl2-dev \
						# libssl-dev \
						# libpq-dev \
						# python-pyinotify \
            # python-renderpm \
            # python-support \
						# libxml2-dev \
						# libxslt1-dev \
						# libldap2-dev \
						# libjpeg-dev \
						# python-dev \
				&& rm -rf /root/.npm /var/lib/apt/lists/* /tmp/* /var/tmp/*

VOLUME ["/var/lib/odoo", "/mnt/extra-addons"]

# Expose Odoo services
EXPOSE 8069 8072

# Set the default config file
ENV ODOO_RC /etc/odoo/odoo.conf

# Set PYTHONPATH
ENV PYTHONPATH /usr/src/odoo

# Set default user when running the container
USER odoo

ENTRYPOINT ["/entrypoint.sh"]
CMD ["odoo"]
