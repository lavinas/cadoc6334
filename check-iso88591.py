def is_iso8859_1(s: str) -> bool:
  """
  Verifica se a string contém apenas caracteres válidos para a codificação ISO-8859-1.

  Args:
    s: A string a ser verificada.

  Returns:
    True se a string for válida em ISO-8859-1, False caso contrário.
  """
  try:
    # Codifica a string para bytes usando ISO-8859-1
    # Se houver um caractere inválido, um UnicodeEncodeError será lançado
    s.encode("iso-8859-1")
    return True
  except UnicodeEncodeError:
    return False

# Exemplo de uso
string_valida = "çã"  # Caracteres acentuados são suportados pelo ISO-8859-1
string_invalida = "Farmácias e Cuidados com a Saúde"  # O emoji não é suportado pelo ISO-8859-1

print(f"'{string_valida}' é ISO-8859-1? {is_iso8859_1(string_valida)}")
print(f"'{string_invalida}' é ISO-8859-1? {is_iso8859_1(string_invalida)}")

# Resultado
# 'Olá, mundo!' é ISO-8859-1? True
# 'Olá, mundo 👋' é ISO-8859-1? False