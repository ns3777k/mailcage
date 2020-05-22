import parseISO from 'date-fns/parseISO';
import formatDistanceToNow from 'date-fns/formatDistanceToNow';

/**
 * Converts "2020-04-18T22:09:46.86556104+03:00" into "8 minutes ago".
 *
 * @param {string} date
 * @returns {string}
 */
export function humanizeDateDistance(date) {
  return `${formatDistanceToNow(parseISO(date))} ago`;
}
