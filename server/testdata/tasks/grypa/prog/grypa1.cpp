#include<iostream>
#define MOD 1000000007
using namespace std;
unsigned long long tab[3][3];//macierz początkowa
unsigned long long wy[3][3]; //macierz wynikowa

void mno(){
    long long pom[3][3];
    pom[0][0]=wy[0][0];
    pom[0][1]=wy[0][1];
    pom[0][2]=wy[0][2];
    pom[1][0]=wy[1][0];
    pom[1][1]=wy[1][1];
    pom[1][2]=wy[1][2];
    pom[2][0]=wy[2][0];
    pom[2][1]=wy[2][1];
    pom[2][2]=wy[2][2];
    
    for (int i=0; i<3; i++)
    {
        for (int j=0; j<3; j++)
        {
            wy[i][j]=0;
            for (int k=0; k<3; k++)
                wy[i][j]+=(pom[i][k]*tab[k][j])%MOD;
            wy[i][j]%=MOD;
        }
    }
}

void kwa(){
    long long pom[3][3];
    pom[0][0]=tab[0][0];
    pom[0][1]=tab[0][1];
    pom[0][2]=tab[0][2];
    pom[1][0]=tab[1][0];
    pom[1][1]=tab[1][1];
    pom[1][2]=tab[1][2];
    pom[2][0]=tab[2][0];
    pom[2][1]=tab[2][1];
    pom[2][2]=tab[2][2];
    
    for (int i=0; i<3; i++)
    {
        for (int j=0; j<3; j++)
        {
            tab[i][j]=0;
            for (int k=0; k<3; k++)
                tab[i][j]+=(pom[i][k]*pom[k][j])%MOD;
            tab[i][j]%=MOD;
        }
    }
}


int main(){
    ios_base::sync_with_stdio(0);
    int n; //liczba dni
    cin>>n;
    int m; //początkowa liczba osób zarażonych
    cin>>m;
    int x,y,z; //współczynniki zaraźliwości
    cin>>x>>y>>z;
    tab[0][0]=0;  //00z
    tab[0][1]=0;  //10y
    tab[0][2]=z;  //01x
    tab[1][0]=1;
    tab[1][1]=0;
    tab[1][2]=y;
    tab[2][0]=0;
    tab[2][1]=1;
    tab[2][2]=x;
    
    wy[0][0]=0;
    wy[0][1]=0;
    wy[0][2]=z;
    wy[1][0]=1;
    wy[1][1]=0;
    wy[1][2]=y;
    wy[2][0]=0;
    wy[2][1]=1;
    wy[2][2]=x;
    n+=1;
    while(n>0)//szybkie potęgowanie
  {

    if (n%2 == 1)
      mno();
 
    kwa();
    n/=2;
  }

  cout<<(wy[2][0]*m)%MOD;
    
}