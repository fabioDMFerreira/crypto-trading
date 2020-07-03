module.exports = {
    "env": {
        "browser": true,
        "es6": true
    },
    "extends": [
        "plugin:react/recommended",
        "airbnb"
    ],
    "globals": {
        "Atomics": "readonly",
        "SharedArrayBuffer": "readonly"
    },
    "parser": "@typescript-eslint/parser",
    "parserOptions": {
        "ecmaFeatures": {
            "jsx": true
        },
        "ecmaVersion": 11,
        "sourceType": "module"
    },
    "plugins": [
        "react",
        "@typescript-eslint",
        'simple-import-sort',
        'jest'
    ],
    "rules": {
        "import/no-unresolved": "off",
        'react/jsx-filename-extension': [2, { extensions: ['.js', '.jsx', '.ts', '.tsx'] }],
        'import/no-extraneous-dependencies': [2, { devDependencies: ['**/*.test.tsx', '**/*.test.ts', '**/setupTests.js', '**/*.stories.tsx'] }],
        'import/extensions': 0,
        'sort-imports': 0,
        "simple-import-sort/sort": "error",
        'no-multiple-empty-lines': "error",
        'no-shadow': 0,
        'max-len': [2, { "code": 150, "ignoreComments": true }],
        "jsx-a11y/anchor-is-valid": "warn",
        "no-plusplus": 1,
        "no-extra-semi": 0,
        "no-param-reassign": 1,
        "prefer-destructuring": 1,
        "react/no-did-update-set-state": 1,
        "jsx-a11y/heading-has-content": 1,
        "import/prefer-default-export": 1,
        "import/no-named-as-default-member": 1,
        "no-underscore-dangle": 1,
        "no-tabs": ["warn", { "allowIndentationTabs": true }],
        "no-continue": 1,
        "sort-imports": "off",
        "import/order": "off",
        "import/first": "off",
        "@typescript-eslint/explicit-function-return-type": "off",
        '@typescript-eslint/no-unused-vars': [2, { args: 'none' }],
        "react/jsx-props-no-spreading": 0
    },
    "settings": {
        'import/parsers': {
            '@typescript-eslint/parser': ['.ts', '.tsx'],
        },
        "import/ignore": [".css$"]
    },
    "env": {
        "jest/globals": true,
        "browser": true
    }
};
