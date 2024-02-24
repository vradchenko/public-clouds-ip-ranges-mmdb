import maxminddb

mmdb_file = 'aws.mmdb'

with maxminddb.open_database(mmdb_file) as reader:
  response = reader.get('13.34.78.160')
  print(response['network_border_group'])
  print(response['region'])
  print(response['service'])
