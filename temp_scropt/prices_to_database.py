import pandas as pd
import os




# Read excel file into dataframe
# usecols argument is a list of integer positions of columns to read in
# pandas will assign column names as 0, 1, 2, etc.
df = pd.read_excel('temp_scropt/prices.xlsx', sheet_name='Виенто', usecols=[1, 2, 8])

print("1 Lines in dataframe:", df.shape[0], "Columns in dataframe:", df.shape[1])
# Drop any row with missing values
# subset argument is a list of column names to check for NaN values
df = df.dropna(subset=[df.columns[0], df.columns[1]])

print("2 Lines in dataframe:", df.shape[0], "Columns in dataframe:", df.shape[1])
# Drop any row where the second column is zero
# .apply() function applies a lambda function to each element of the column
# .notnull() and .notna() check for non-missing values
# df = df[(df.iloc[:, 1].apply(lambda x: x != 0 if not pd.isna(x) and not pd.isnull(x) else True))]

# Drop any row where the second column is not a numeric value
# .to_numeric() with errors='coerce' will convert non-numeric values to NaN
# .notnull() will drop any rows with NaN values
print("First 10 lines of dataframe:")
print(df.head(10))
# Drop any row where the first column is non-digit
df = df[pd.to_numeric(df.iloc[:, 0], errors='coerce').notnull()]


print("3 Lines in dataframe:", df.shape[0], "Columns in dataframe:", df.shape[1])
# Reorder the columns
# just pass a list of column names in the desired order
df.iloc[:, 0] = df.iloc[:, 0].astype(str)
df.iloc[:, 1] = df.iloc[:, 1].apply(lambda x: str(x).replace(',', ''))
df = df[[df.columns[1], df.columns[2], df.columns[0]]]

# Write dataframe to excel file
# index=False argument will not write row index values
df.to_csv('temp_scropt/prices.csv', index=False)

