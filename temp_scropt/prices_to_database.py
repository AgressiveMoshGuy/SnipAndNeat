import pandas as pd


df = pd.read_excel('temp_scropt/prices.xlsx', sheet_name='Виенто', usecols=[1, 2, 8])
df = df.dropna(subset=[df.columns[0], df.columns[1]])
df = df[(pd.to_numeric(df.iloc[:, 1], errors='coerce').notnull())]
df = df[pd.to_numeric(df.iloc[:, 0], errors='coerce').notnull()]
df.iloc[:, 0] = df.iloc[:, 0].astype(str)
df = df[[df.columns[1], df.columns[2], df.columns[0]]]
df.to_excel('out.xlsx', index=False)

