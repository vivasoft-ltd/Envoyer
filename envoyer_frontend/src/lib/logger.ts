import winston, { format } from 'winston';

const { colorize, combine, timestamp, printf } = format;

const logFormat = printf(({ level, message, timestamp }) => {
  return `|${timestamp}| ${level}: ${message}`;
});

const transports: winston.transport[] = [new winston.transports.Console()];

const logger = winston.createLogger({
  level: 'info',
  format: combine(colorize(), timestamp(), logFormat),
  transports,
});

export default logger;
