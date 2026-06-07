# User Guide

This guide explains how to interact with the system once the environment is configured.

## 1. Product Management
- **Add product**: Copy the URL of a product from a supported store. Paste it into the input field and click "Add".
- **Search products**: Typing the name of an existing product will display the product or products sharing those names on the main page.
- **Update**: Each product has a manual update button if you need to force a real-time price check.
- **Creating Notifications**: On the product page, if no price is entered, the historical lowest price of the product (since its creation) will be set as the alarm. If a value is entered, it will only trigger when the price drops below that value.

## 2. Price Tracking
- **Target prices**: You can define a maximum price to receive alerts. The system compares the current price with your target.
- **History**: In the product detail view, a price history chart is rendered based on the history saved in the database.

## 3. Troubleshooting

| Issue | Likely cause | Solution |
| :--- | :--- | :--- |
| Blank product page | Blocked by the website or empty HTML | Access the product link directly in your browser to verify it works correctly. If it does, wait a bit with the website open and update the product again |
| A specific website doesn't work (Missing data or doesn't load) | It is not implemented in the scraper, so the HTML tags used by the website are not found | Get in touch and the website will be included in future versions |

## 4. Workflow for new contributors
1. **Setup:** Follow the steps in the `README.md`.
2. **Develop:** Create a `feature/` or `fix/` branch.
3. **Test:** Run tests
4. **Pull Request:** Open a PR to `main` after validating the changes.