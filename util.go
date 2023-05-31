package cfg

// O representa um objeto JSON, usado para criação de configurações padrão
type O map[string]interface{}

// Merge faz a mesclagem do objeto atual com o outro sem manter referencia
func (o O) Merge(other O) {
   if other == nil {
      return
   }
   for key, value := range other {
      switch v := value.(type) {
      case O:
         vl := O{}
         vl.Merge(v)
         o[key] = vl
         break
      case []string:
         o[key] = v[0:]
         break
      case []O:
         var vl []O
         for _, it := range v {
            vli := O{}
            vli.Merge(it)
            vl = append(vl, vli)
         }
         o[key] = vl
         break
      default:
         o[key] = value
      }
   }
}

type DefaultConfigFn func(extra ...O) O

// CreateDefaultConfigFn simplifica a criação de configurações genéricas
func CreateDefaultConfigFn(defaultValue O) DefaultConfigFn {
   return func(extra ...O) O {
      config := O{}
      config.Merge(defaultValue)

      for _, o2 := range extra {
         config.Merge(o2)
      }
      return config
   }
}
