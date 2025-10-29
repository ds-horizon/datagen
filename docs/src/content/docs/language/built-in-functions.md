---
title: Built-in Functions
description: Comprehensive reference for datagen's built-in functions
---

datagen provides a comprehensive set of built-in functions for common data generation needs. These functions are optimized for performance and provide realistic data patterns.

## Function Reference Summary

| Category            | Function                                  | Return Type              | Description                                    |
|---------------------|-------------------------------------------|--------------------------|------------------------------------------------|
| **Personal**        | `Name()`                                  | `string`                 | Random full name                               |
|                     | `Email()`                                 | `string`                 | Random email address                           |
|                     | `Phone()`                                 | `string`                 | Random phone number                            |
|                     | `FirstName()`                             | `string`                 | Random first name                              |
|                     | `LastName()`                              | `string`                 | Random last name                               |
|                     | `MiddleName()`                            | `string`                 | Random middle name                             |
|                     | `NamePrefix()`                            | `string`                 | Random name prefix (e.g., Mr., Dr.)            |
|                     | `NameSuffix()`                            | `string`                 | Random name suffix (e.g., Jr., Sr.)            |
|                     | `Gender()`                                | `string`                 | Random gender                                  |
|                     | `Username()`                              | `string`                 | Random username                                |
|                     | `SSN()`                                   | `string`                 | Random SSN                                     |
|                     | `PhoneFormatted()`                        | `string`                 | Random formatted phone                         |
| **Text**            | `Text(int)`                         | `string`                 | Random text of specified length                |
|                     | `Sentence(int)`                           | `string`                 | Random sentence                                |
| **Numbers**         | `IntBetween(int, int)`              | `int`                    | Random integer in range                        |
|                     | `FloatBetween(float64, float64)`    | `float64`                | Random float64 in range                        |
|                     | `Float32Between(float32, float32)`  | `float32`                | Random float32 in range                        |
|                     | `Boolean()`                         | `bool`                   | Random boolean                                 |
| **Dates**           | `DateBetween(time.Time, time.Time)` | `time.Time`              | Random date in range                           |
|                     | `DateBetweenStr(string, string)`    | `string`                 | Random date string                             |
|                     | `Date()`                                  | `time.Time`              | Random date                                    |
|                     | `FutureDate()`                            | `time.Time`              | Random future date                             |
|                     | `PastDate()`                              | `time.Time`              | Random past date                               |
|                     | `NanoSecond()`                            | `int`                    | Random nanosecond                              |
|                     | `Second()`                                | `int`                    | Random second                                  |
|                     | `Minute()`                                | `int`                    | Random minute                                  |
|                     | `Hour()`                                  | `int`                    | Random hour                                    |
|                     | `Month()`                                 | `int`                    | Random month number                            |
|                     | `MonthString()`                           | `string`                 | Random month name                              |
|                     | `Day()`                                   | `int`                    | Random day of month                            |
|                     | `WeekDay()`                               | `string`                 | Random weekday name                            |
|                     | `Year()`                                  | `int`                    | Random year                                    |
|                     | `TimeZone()`                              | `string`                 | Random timezone                                |
|                     | `TimeZoneAbv()`                           | `string`                 | Random timezone abbr                           |
|                     | `TimeZoneFull()`                          | `string`                 | Random timezone full name                      |
|                     | `TimeZoneOffset()`                        | `float32`                | Random timezone offset                         |
|                     | `TimeZoneRegion()`                        | `string`                 | Random timezone region                         |
| **Selection**       | `RandomFrom(T...)`                        | `T`                      | Random selection from values                   |
| **IDs**             | `UUID()`                                  | `string`                 | Random UUID                                    |
| **Business**        | `Company()`                         | `string`                 | Random company name                            |
|                     | `Url()`                             | `string`                 | Random URL                                     |
|                     | `HexColor()`                        | `string`                 | Random hex color                               |
| **Conversion**      | `ToJSON(interface{})`                     | `string`                 | Convert to JSON string                         |
| **Security**        | `Password(PasswordOptions)`               | `string`                 | Random password with options                   |
| **Address**         | `AddressInfo()`                           | `string`                 | Random full address string                     |
|                     | `AddressFull()`                           | `string`                 | Random full address (street, city, state, zip) |
|                     | `Street()`                                | `string`                 | Random street                                  |
|                     | `StreetName()`                            | `string`                 | Random street name                             |
|                     | `StreetNumber()`                          | `string`                 | Random street number                           |
|                     | `StreetPrefix()`                          | `string`                 | Random street prefix                           |
|                     | `StreetSuffix()`                          | `string`                 | Random street suffix                           |
|                     | `City()`                                  | `string`                 | Random city                                    |
|                     | `State()`                                 | `string`                 | Random state                                   |
|                     | `StateAbr()`                              | `string`                 | Random state abbreviation                      |
|                     | `Zip()`                                   | `string`                 | Random ZIP code                                |
|                     | `Country()`                               | `string`                 | Random country                                 |
|                     | `CountryAbr()`                            | `string`                 | Random country abbreviation                    |
| **Geography**       | `Latitude()`                              | `float64`                | Random latitude                                |
|                     | `Longitude()`                             | `float64`                | Random longitude                               |
|                     | `LatitudeInRange(float64, float64)`       | `(float64, error)`       | Latitude within range                          |
|                     | `LongitudeInRange(float64, float64)`      | `(float64, error)`       | Longitude within range                         |
| **Internet**        | `DomainName()`                            | `string`                 | Random domain name                             |
|                     | `DomainSuffix()`                          | `string`                 | Random domain suffix                           |
|                     | `IPv4Address()`                           | `string`                 | Random IPv4 address                            |
|                     | `IPv6Address()`                           | `string`                 | Random IPv6 address                            |
|                     | `MacAddress()`                            | `string`                 | Random MAC address                             |
|                     | `HTTPMethod()`                            | `string`                 | Random HTTP method                             |
|                     | `HTTPStatusCode()`                        | `int`                    | Random HTTP status code                        |
|                     | `HTTPStatusCodeSimple()`                  | `int`                    | Common HTTP status code                        |
|                     | `HTTPVersion()`                           | `string`                 | Random HTTP version                            |
|                     | `UserAgent()`                             | `string`                 | Random user agent                              |
|                     | `ChromeUserAgent()`                       | `string`                 | Chrome user agent                              |
|                     | `FirefoxUserAgent()`                      | `string`                 | Firefox user agent                             |
|                     | `SafariUserAgent()`                       | `string`                 | Safari user agent                              |
|                     | `OperaUserAgent()`                        | `string`                 | Opera user agent                               |
| **Numeric Types**   | `Int8()`                                  | `int8`                   | Random int8                                    |
|                     | `Int16()`                                 | `int16`                  | Random int16                                   |
|                     | `Int32()`                                 | `int32`                  | Random int32                                   |
|                     | `Int64()`                                 | `int64`                  | Random int64                                   |
|                     | `Uint8()`                                 | `uint8`                  | Random uint8                                   |
|                     | `Uint16()`                                | `uint16`                 | Random uint16                                  |
|                     | `Uint32()`                                | `uint32`                 | Random uint32                                  |
|                     | `Uint64()`                                | `uint64`                 | Random uint64                                  |
|                     | `Float32()`                               | `float32`                | Random float32                                 |
|                     | `Float64()`                               | `float64`                | Random float64                                 |
|                     | `Int()`                                   | `int`                    | Random int                                     |
|                     | `IntN(int)`                               | `int`                    | Random int within N                            |
|                     | `Uint()`                                  | `uint`                   | Random uint                                    |
|                     | `UintN(uint)`                             | `uint`                   | Random uint within N                           |
|                     | `UintRange(uint, uint)`                   | `uint`                   | Random uint in range                           |
|                     | `HexUint(int)`                            | `string`                 | Random hex number string                       |
| **Colors**          | `Color()`                                 | `string`                 | Random color name                              |
|                     | `SafeColor()`                             | `string`                 | Random safe color name                         |
|                     | `RGBColor()`                              | `[]int`                  | Random RGB triplet                             |
|                     | `NiceColors()`                            | `[]string`               | Palette of nice colors                         |
| **Files & Media**   | `FileExtension()`                         | `string`                 | Random file extension                          |
|                     | `FileMimeType()`                          | `string`                 | Random file MIME type                          |
|                     | `Image(int, int)`                         | `*image.RGBA`            | Random image RGBA                              |
|                     | `ImageJpeg(int, int)`                     | `[]byte`                 | Random JPEG image bytes                        |
|                     | `ImagePng(int, int)`                      | `[]byte`                 | Random PNG image bytes                         |
|                     | `SvgString(SVGOpts)`                      | `string`                 | Random SVG string                              |
| **Finance**         | `CreditCardNumber()`                      | `string`                 | Random credit card number                      |
|                     | `CreditCardType()`                        | `string`                 | Random credit card type                        |
|                     | `CreditCardExp()`                         | `string`                 | Random credit card expiry                      |
|                     | `CreditCardCvv()`                         | `string`                 | Random credit card CVV                         |
|                     | `Currency()`                              | `string`                 | Random currency code                           |
|                     | `CurrencyLong()`                          | `string`                 | Random currency name                           |
| **Jobs**            | `JobTitle()`                              | `string`                 | Random job title                               |
|                     | `JobDescriptor()`                         | `string`                 | Random job descriptor                          |
|                     | `JobLevel()`                              | `string`                 | Random job level                               |
| **Locales & Lang**  | `Language()`                              | `string`                 | Random language                                |
|                     | `LanguageBCP()`                           | `string`                 | Random BCP-47 language tag                     |
|                     | `ProgrammingLanguage()`                   | `string`                 | Random programming language                    |
| **Emojis**          | `Emoji()`                                 | `string`                 | Random emoji                                   |
| **Company**         | `CompanySuffix()`                         | `string`                 | Random company suffix                          |
| **HTTP/Logging**    | `LogLevel()`                              | `string`                 | Random log level                               |
| **Text Utils**      | `Digit()`                                 | `string`                 | Random digit                                   |
|                     | `DigitN(uint)`                            | `string`                 | Random N digits                                |
|                     | `Letter()`                                | `string`                 | Random letter                                  |
|                     | `Numerify(string)`                        | `string`                 | Replace numerals in pattern                    |
|                     | `Word()`                                  | `string`                 | Random word                                    |
|                     | `Paragraph(ParagraphOptions)`             | `string`                 | Random lorem paragraph                         |
|                     | `Regex(string)`                           | `string`                 | Random string by regex                         |
| **Data Structures** | `Map()`                                   | `map[string]interface{}` | Random map                                     |
|                     | `Slice(interface{})`                      | `void`                   | Randomize slice in-place                       |
|                     | `Struct(interface{})`                     | `error`                  | Fill struct with random data                   |
| **Utilities**       | `InputName()`                             | `string`                 | Random HTML input name                         |
|                     | `Unit()`                                  | `string`                 | Random measurement unit                        |
|                     | `Weighted([]any, []float32)`              | `(any, error)`           | Weighted random choice                         |
|                     | `ShuffleInts([]int)`                      | `void`                   | Shuffle int slice                              |
|                     | `ShuffleStrings([]string)`                | `void`                   | Shuffle string slice                           |
|                     | `ShuffleAnySlice(any)`                    | `void`                   | Shuffle any slice                              |
| **Errors**          | `Error()`                                 | `error`                  | Random generic error                           |
|                     | `ErrorObject()`                           | `error`                  | Random error object                            |
|                     | `ErrorDatabase()`                         | `error`                  | Random database error                          |
|                     | `ErrorGRPC()`                             | `error`                  | Random gRPC error                              |
|                     | `ErrorHTTP()`                             | `error`                  | Random HTTP error                              |
|                     | `ErrorRuntime()`                          | `error`                  | Random runtime error                           |
|                     | `ErrorValidation()`                       | `error`                  | Random validation error                        |
|                     | `ErrorHTTPClient()`                       | `error`                  | Random HTTP client error                       |
|                     | `ErrorHTTPServer()`                       | `error`                  | Random HTTP server error                       |
| **Misc**            | `Fruit()`                                 | `string`                 | Random fruit name                              |
|                     | `Vegetable()`                             | `string`                 | Random vegetable name                          |
| **Text Utils**      | `MarkdownSimple()`                        | `(string, error)`        | Random markdown string                         |
