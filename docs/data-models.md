# Data Models

## User
Represents the user
| Field | Tipe | Description |
| :--- | :--- | :--- |
| `UserID` | uint | Unic Identifier (PK) |
| `Username` | string | User´s nickname |
| `Email` | string | User´s associated email |
| `Password` | string | User´s associated password |
| `UserType` | string | Type of user |
| `RegisterAt` | time.Time | Date the user was registered |
| `SavedProducts` | pq.Int64Array | Array with followed products´ ID |

### Product
Represents the scrapped product on the platform
| Field | Tipe | Description |
| :--- | :--- | :--- |
| `ProductID` | uint | Unic Identifier (PK) |
| `Name` | string | Product´s name |
| `SourceURL` | string | Original store URL (unique index) |
| `LastPrice` | float64 | Last recorded price |
| `LowestPrice` | float64 | Lowest historical price recorded |
| `CreatedBy` | uint | FK of the user that created the product |
| `CreateAt` | time.Time | Date when the product was created |
| `UpdatedAt` | time.Time | Date when the product was last updated |


### Tracking
Defines the relationship between users and products for notifications
| Field | Tipe | Description |
| :--- | :--- | :--- |
| `TrackingID` | uint | Unique identifier (PK) |
| `UserID` | uint | FK of the user following the product |
| `ProductID` | uint | FK of the followed product |
| `NotifyPriceChanges`| bool | Flag to enable notifications |
| `NotifyTargetPrice`| float64 | Value for the notification |
| `TrackingStartDate` | time.Time | Date when the product was started tracking |


### PriceHistory
Stores historical price records
| Field | Tipe | Description |
| :--- | :--- | :--- |
| `PriceHistoryID`| uint | Unique identifier (PK) |
| `ProductID` | uint | FK referencing `Product` |
| `Price` | float64 | Price on the recorded date |
| `RegisterDate` | time.Time | Registration timestamp |

### Notifications
Manages alerts sent to users regarding price updates or status changes
| Field | Tipe | Description |
| :--- | :--- | :--- |
| `NotificationID` | uint | Unic Identifier (PK) |
| `UserID` | uint | FK referencing the recipient `User` |
| `ProductID` | uint | FK referencing the related `Product` |
| `Type` | string | Notification category |
| `Title` | string | Brief notification headline |
| `Description` | string | Detailed content of the alert |
| `IsRead` | bool | Status flag |
| `CreatedAt` | time.Time | Timestamp of notification creation |
