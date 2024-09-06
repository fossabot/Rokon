PORTNAME=	     rokon
PORTVERSION=     1.0.0
DISTNAME=	${PORTNAME}${PORTVERSION}
CATEGORIES=	     utility
MASTER_SITES=    https://github.com/BrycensRanch/Rokon/archive/master.tar.gz
MAINTAINER=      brycengranville@outlook.com
COMMENT=         Control your Roku device with your desktop!
WWW=             https://github.com/BrycensRanch/Rokon
LICENSE=         AGPL-3.0-or-later
USES=		     tar:tgz compiler:c11 gmake go
# Build dependencies: Go, GTK4, git, gcc
BUILD_DEPENDS=	${LOCALBASE}/bin/go:lang/go \
		${LOCALBASE}/bin/gtk4:devel/gtk4 \
		${LOCALBASE}/bin/git:devel/git \
		${LOCALBASE}/bin/gcc:lang/gcc

# Runtime dependency: GTK4
RUN_DEPENDS=	${LOCALBASE}/bin/gtk4:devel/gtk4

# Logic to append Git commit SHA to the version
GIT_SHA!=	git rev-parse --short HEAD
VERSION_WITH_SHA!= echo ${DISTVERSION} | sed -E "s/(.*)-[a-f0-9]{7}/\1/" && echo ${DISTVERSION} | grep -q "-${GIT_SHA}" || echo "${DISTVERSION}-${GIT_SHA}"


GO_PKGNAME=	Rokon

PLIST_FILES=	bin/rokon \
		share/applications/io.github.brycensranch.Rokon.desktop \
		share/icons/hicolor/48x48/apps/io.github.brycensranch.Rokon.png \
		share/icons/hicolor/256x256/apps/io.github.brycensranch.Rokon.png \
		share/icons/hicolor/scalable/apps/io.github.brycensranch.Rokon.svg \
		share/metainfo/io.github.brycensranch.Rokon.metainfo.xml \
		%%DOCSDIR%%/*.md

OPTIONS_DEFINE=	DOCS

post-install-DOCS-on:
	@${MKDIR} ${STAGEDIR}${DOCSDIR}
	${INSTALL_DATA} ${WRKSRC}/*.md ${STAGEDIR}${DOCSDIR}

.include <bsd.port.mk>
