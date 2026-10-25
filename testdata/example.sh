#!/bin/sh
#
# This program writes words.
#
REPO=~bradley/.m2/repository
TOOLS_VERSION=2.15.0-SNAPSHOT
PICOCLI_VERSION=4.7.7
COMMONS_VERSION=4.0.0
SLF4J_VERSION=2.0.18
LOGBACK_VERSION=1.5.32
PGM_CLASSES=target/classes

cp=$PGM_CLASSES
cp=$cp:${REPO}/org/natuna/commons/commons-core/${COMMONS_VERSION}/commons-core-${COMMONS_VERSION}.jar
cp=$cp:${REPO}/org/natuna/commons/commons-picocli/${COMMONS_VERSION}/commons-picocli-${COMMONS_VERSION}.jar
cp=$cp:${REPO}/info/picocli/picocli/${PICOCLI_VERSION}/picocli-${PICOCLI_VERSION}.jar
cp=$cp:${REPO}/org/slf4j/slf4j-api/${SLF4J_VERSION}/slf4j-api-${SLF4J_VERSION}.jar
cp=$cp:${REPO}/ch/qos/logback/logback-classic/${LOGBACK_VERSION}/logback-classic-${LOGBACK_VERSION}.jar
cp=$cp:${REPO}/ch/qos/logback/logback-core/${LOGBACK_VERSION}/logback-core-${LOGBACK_VERSION}.jar

java -cp "$cp" org.natuna.tools.Words "$@"
