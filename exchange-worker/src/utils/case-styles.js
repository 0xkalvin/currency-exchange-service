function toSnakeCase(key) {
  return key.replace(/([A-Z])/g, '_$1').toLowerCase();
}

const toCamel = (s) => s.replace(/([-_][a-z])/ig, ($1) => $1.toUpperCase()
  .replace('-', '')
  .replace('_', ''));

const isObject = (o) => o === Object(o) && !Array.isArray(o) && typeof o !== 'function';

function objectToSnakeCase(object = {}) {
  const newObject = {};

  Object.keys(object).forEach((key) => {
    const parsedKey = toSnakeCase(key);

    newObject[parsedKey] = object[key];
  });

  return newObject;
}

const objectToCammelCase = function (o) {
  if (isObject(o)) {
    const n = {};

    Object.keys(o)
      .forEach((k) => {
        n[toCamel(k)] = objectToCammelCase(o[k]);
      });

    return n;
  } if (Array.isArray(o)) {
    return o.map((i) => objectToCammelCase(i));
  }

  return o;
};

module.exports = {
  objectToSnakeCase,
  objectToCammelCase,
  toSnakeCase,
};
