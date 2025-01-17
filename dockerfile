FROM golang:1.15-buster


ENV GO111MODULE=on \
   CGO_ENABLED=1 \
   CGO_CFLAGS_ALLOW='-Xpreprocessor'\
   GOOS=linux

WORKDIR /lecture-parser


ENV DEBIAN_FRONTEND noninteractive
ENV AWS_SDK_LOAD_CONFIG true

# RUN apt-get update \
#    && apt-get install -y \
#    wget build-essential \
#    pkg-config \
#    --no-install-recommends \
#    && apt-get -q -y install \
#    libjpeg-dev \
#    libpng-dev \
#    ffmpeg\
#    libtiff-dev \
#    libgif-dev \
#    libx11-dev \
#    fontconfig fontconfig-config libfontconfig1-dev \
#    ghostscript gsfonts gsfonts-x11 \
#    libfreetype6-dev \
#    --no-install-recommends \
#    && rm -rf /var/lib/apt/lists/*g


# use latest version, don't build manually. Comment out the gs rights="none" line in gsconfig. Make sure all the required delegates are installed.
# Path instructions for gocv is outdated on website. $gopath/pkg/mod/gocv
# yay -S opencv hdf5 vtk  
# gocv needs hdf5 and vtk to work. Building doesn't seem to work well.
ENV IMAGEMAGICK_VERSION=7.0.8-11

RUN cd && \
   wget https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz && \
   tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz && \
   cd ImageMagick* && \
   ./configure \
   --without-magick-plus-plus \
   --without-perl \
   --disable-openmp \
   --with-gvc=no \
   --disable-docs && \
   make -j$(nproc) && make install && \
   ldconfig /usr/local/lib

# COPY go.mod .
# COPY go.sum .

COPY . .

RUN go mod download

RUN go build -o lecture-parser ./cmd/lecture-parser

WORKDIR /dist
RUN cp /lecture-parser/lecture-parser .
RUN cp -r /lecture-parser/altInput .

CMD [ "/dist/lecture-parser" ]

